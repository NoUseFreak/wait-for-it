package main

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

type Config struct {

	Services []ServiceConfig
}

func (c *Config) ParseConfig(content []byte) error {

	m := make(map[string]map[string]map[string]string)
	err := yaml.Unmarshal(content, &m)

	if _, ok := m["services"]; !ok {
		log.Fatal("no services set")
	}

	for name, value := range(m["services"]) {
		c.Services = append(c.Services, ServiceConfig{
			Name: name,
			Type: value["type"],
			Settings:value,
		})
	}

	return err
}

func NewConfig(path string) (Config, error) {
	config := Config{}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	config.ParseConfig(content)

	return config, nil
}

type ServiceConfig struct {
	Name string
	Type string
	Settings map[string]string
}