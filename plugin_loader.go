package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

type PluginLoader struct {
	location          string
	latestReleasePath string
}

func NewPluginLoader(location string) (PluginLoader, error) {
	pluginLoader := PluginLoader{
		location: location,
	}

	return pluginLoader, nil
}

func (pl *PluginLoader) ensureDirectory() {
	if _, err := os.Stat(pl.location); os.IsNotExist(err) {
		os.MkdirAll(pl.location, os.ModePerm)
	}
}

func (pl *PluginLoader) LoadAll(configs []ServiceConfig, noCache bool) {
	cliUi.Title("Initializing plugins...")
	if noCache {
		pl.CleanUp()
	}
	pl.ensureDirectory()

	plugins := map[string]int{}
	for _, service := range configs {
		plugins[service.Type]++
	}

	var wg sync.WaitGroup
	wg.Add(len(plugins))
	for name, _ := range plugins {
		go func(name string) {
			defer wg.Done()
			err := pl.LoadPlugin(name)
			if err != nil {
				cliUi.Error("Failed to load " + name)
				os.Exit(1)
			}
		}(name)
	}

	wg.Wait()
}

func (pl *PluginLoader) LoadPlugin(name string) error {
	targetPath := pl.location + "/" + name
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		cliUi.Output(fmt.Sprintf(" - Downloading plugin for %s...", name))
		pluginUrl, err := pl.findPluginUrl(name)
		if err != nil {
			return err
		}

		return pl.downloadPlugin(pluginUrl, name)
	} else {
		cliUi.Output(fmt.Sprintf(" - Downloading plugin for %s... (cached)", name))
	}
	return nil
}

func (pl *PluginLoader) findPluginUrl(name string) (string, error) {
	url, err := pl.findLatestReleasePath()
	if err != nil {
		return "", err
	}

	return url + "/" + runtime.GOOS + "_" + name, nil
}

func (pl *PluginLoader) findLatestReleasePath() (string, error) {
	if pl.latestReleasePath != "" {
		return pl.latestReleasePath, nil
	}

	releaseUrl := "https://api.github.com/repos/NoUseFreak/wait-for-it/releases/latest"
	client := http.Client{Timeout: time.Second * 2}
	req, err := http.NewRequest(http.MethodGet, releaseUrl, nil)
	if err != nil {
		return "", err
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		cliUi.Error(getErr.Error())
		return "", getErr
	}

	data := map[string]interface{}{}
	dec := json.NewDecoder(res.Body)
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	value, err := jq.String("assets", "0", "browser_download_url")
	if err != nil {
		return "", nil
	}
	assetUrl, err := url.Parse(value)
	if err != nil {
		return "", nil
	}

	pl.latestReleasePath = strings.Replace(value, assetUrl.Path, path.Dir(assetUrl.Path), 1)

	return pl.latestReleasePath, nil
}

func (pl *PluginLoader) downloadPlugin(url string, name string) error {
	target := pl.location + "/" + name
	output, err := os.Create(target)
	if err != nil {
		fmt.Println("Error while creating", target, "-", err)
		return err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		os.Remove(target)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("Error while downloading", url, "-", err)
		os.Remove(target)
		os.Exit(1)
	}

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		os.Remove(target)
		os.Exit(1)
		return err
	}
	os.Chmod(target, 0755)

	return nil
}
func (pl *PluginLoader) CleanUp() error {
	return os.RemoveAll(pl.location)
}
