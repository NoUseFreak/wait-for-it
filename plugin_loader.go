package main

import (
	"fmt"
	"os"
)

type PluginLoader struct {
	location string
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
	for _, service := range configs {
		pl.LoadPlugin(service.Type)
	}
}

func (pl *PluginLoader) LoadPlugin(name string) error {
	path := pl.location + "/" + name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cliUi.Output(fmt.Sprintf(" - Downloading plugin for %s...", name))
	} else {
		cliUi.Output(fmt.Sprintf(" - Downloading plugin for %s... (cached)", name))
	}
	return nil
}

func (pl *PluginLoader) CleanUp() error {
	return os.RemoveAll(pl.location)
}
