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

	m := make(map[string]map[string]map[string]interface{})
	if err := yaml.Unmarshal(content, &m); err != nil {
		return err
	}

	if _, ok := m["services"]; !ok {
		return fmt.Errorf("no service found")
	}

	for name, value := range m["services"] {
		sc, err := NewServiceConfig(name, value, *c)
		if err != nil {
			return err
		}
		c.Services = append(c.Services, sc)
	}

	return nil
}

type ServiceConfig struct {
	Name     string
	Type     string
	Timeout  time.Duration
	Settings map[string]string
}

func NewServiceConfig(name string, settings map[string]interface{}, defaults Config) (ServiceConfig, error) {
	sc := ServiceConfig{}

	if _, ok := settings["plugin"]; !ok {
		return sc, fmt.Errorf("no service found")
	}
	if _, ok := settings["parameters"]; !ok {
		return sc, fmt.Errorf("no service found")
	}

	sc.Name = name
	sc.Type = settings["plugin"].(string)
	sc.Timeout = defaults.DefaultTimeout

	sc.Settings = make(map[string]string)
	for key, value := range settings["parameters"].(map[interface{}]interface{}) {
		sc.Settings[key.(string)] = value.(string)
	}

	if val, ok := settings["timeout"]; ok {
		valInt, _ := strconv.Atoi(val.(string))
		sc.Timeout = time.Duration(valInt) * time.Second
	}

	delete(sc.Settings, "type")
	delete(sc.Settings, "timeout")

	return sc, nil
}
