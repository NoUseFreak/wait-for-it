package main

import (
	"fmt"
	"os"
	"net/http"
	"time"
	"encoding/json"
	"github.com/jmoiron/jsonq"
	"path"
	"io"
	"net/url"
	"strings"
)

type PluginLoader struct {
	location string
	latestReleasePath string
}

func NewPluginLoader(location string) (PluginLoader, error) {
	pluginLoader := PluginLoader{
		location: location,
	}
	pluginLoader.ensureDirectory()

	return pluginLoader, nil
}

func (pl *PluginLoader) ensureDirectory() {
	if _, err := os.Stat(pl.location); os.IsNotExist(err) {
		os.MkdirAll(pl.location, os.ModePerm)
	}
}

func (pl *PluginLoader) LoadAll(configs []ServiceConfig) {
	cliUi.Title("Initializing plugins...")

	plugins := map[string]int{}

	for _, service := range configs {
		plugins[service.Type]++
	}

	fmt.Println(plugins)

	for name, _ := range plugins {
		pl.LoadPlugin(name)
	}
}

func (pl *PluginLoader) LoadPlugin(name string) error {
	targetPath := pl.location + "/" + name
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		cliUi.Output(fmt.Sprintf(" - Downloading plugin for %s...", name))
		pluginUrl, _ := pl.findPluginUrl(name)
		pl.downloadPlugin(pluginUrl, name)
	} else {
		cliUi.Output(fmt.Sprintf(" - Downloading plugin for %s... (cached)", name))
	}
	return nil
}

func (pl *PluginLoader) findPluginUrl(name string) (string, error) {
	url, _ := pl.findLatestReleasePath()

	return url + "/" + name, nil
}

func (pl *PluginLoader) findLatestReleasePath() (string, error) {
	if pl.latestReleasePath != "" {
		return pl.latestReleasePath, nil
	}

	releaseUrl := "https://api.github.com/repos/NoUseFreak/wait-for-it/releases/latest"
	client := http.Client{Timeout: time.Second * 2}
	req, _ := http.NewRequest(http.MethodGet, releaseUrl, nil)
	res, getErr := client.Do(req)
	if getErr != nil {
		cliUi.Error(getErr.Error())
		return "", getErr
	}

	data := map[string]interface{}{}
	dec := json.NewDecoder(res.Body)
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	value, _ := jq.String("assets", "0", "browser_download_url")
	assetUrl, _ := url.Parse(value)

	pl.latestReleasePath = strings.Replace(value, assetUrl.Path, path.Dir(assetUrl.Path), 1)

	return pl.latestReleasePath, nil
}

func (pl *PluginLoader) downloadPlugin(url string, name string) {
	target := pl.location + "/" + name
	output, err := os.Create(target)
	if err != nil {
		fmt.Println("Error while creating", target, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	os.Chmod(target, 0755)

	fmt.Println(n, "bytes downloaded.")
}
func (pl *PluginLoader) CleanUp() error {
	return os.RemoveAll(pl.location)
}
