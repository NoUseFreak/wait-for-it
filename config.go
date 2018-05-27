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
	if err := config.ParseConfig(content); err != nil {
		return config, err
	}

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
	delete(sc.Settings, "plugin")

	sc.Settings = make(map[string]string)
	for key, value := range settings["parameters"].(map[interface{}]interface{}) {
		switch t := value.(type) {
		case string:
			sc.Settings[key.(string)] = value.(string)
		case int:
			sc.Settings[key.(string)] = strconv.Itoa(value.(int))
		default:
			return sc, fmt.Errorf("can't handle value of type %T", t)
		}
	}

	sc.Timeout = defaults.DefaultTimeout
	if val, ok := settings["timeout"]; ok {
		switch t := val.(type) {
		case string:
			intVal, _ := strconv.Atoi(val.(string))
			sc.Timeout = time.Duration(intVal) * time.Second
		case int:
			sc.Timeout = time.Duration(val.(int)) * time.Second
		default:
			return sc, fmt.Errorf("unknown type for timeout %T", t)
		}
	}
	delete(sc.Settings, "timeout")

	return sc, nil
}
