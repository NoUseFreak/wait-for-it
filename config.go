package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"time"
	"fmt"
)

type Config struct {
	DefaultTimeout time.Duration
	Services       []ServiceConfig
}

func NewConfig(path string) (Config, error) {
	config := Config{}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	config.SetDefaults()
	config.ParseConfig(content)

	return config, nil
}

func (c *Config) SetDefaults() {
	c.DefaultTimeout = 30 * time.Second
}

func (c *Config) ParseConfig(content []byte) error {

	m := make(map[string]map[string]map[string]string)
	err := yaml.Unmarshal(content, &m)

	if _, ok := m["services"]; !ok {
		return fmt.Errorf("no service found")
	}

	for name, value := range m["services"] {
		sc, _ := NewServiceConfig(name, value, *c)
		c.Services = append(c.Services, sc)
	}

	return err
}

type ServiceConfig struct {
	Name     string
	Type     string
	Timeout  time.Duration
	Settings map[string]string
}

func NewServiceConfig(name string, settings map[string]string, defaults Config) (ServiceConfig, error) {
	sc := ServiceConfig{
		Name:     name,
		Type:     settings["type"],
		Timeout:  defaults.DefaultTimeout,
		Settings: settings,
	}

	if val, ok := settings["timeout"]; ok {
		valInt, _ := strconv.Atoi(val)
		sc.Timeout = time.Duration(valInt) * time.Second
	}

	delete(sc.Settings, "type")
	delete(sc.Settings, "timeout")

	return sc, nil
}
