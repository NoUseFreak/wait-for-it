package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestNewConfig_Existing(t *testing.T) {
	_, err := NewConfig("wait-for-it.yml")

	assert.Nil(t, err)
}

func TestNewConfig_NonExistentFile(t *testing.T) {
	_, err := NewConfig("this-is-not-here.yml")

	assert.Error(t, err)
}

func TestConfig_SetDefaults(t *testing.T) {
	config := Config{}
	config.SetDefaults()

	assert.NotEmpty(t, config.DefaultTimeout)
}

func TestConfig_ParseConfig_Success(t *testing.T) {
	var data = `
services:
  mysql_check:
    plugin: mysql
`
	config := Config{}
	err := config.ParseConfig([]byte(data))

	assert.Nil(t, err)
	assert.NotEmpty(t, config.Services)
}

func TestConfig_ParseConfig_FailOnNoServices(t *testing.T) {
	var data = `
random:
  mysql_check:
    plugin: mysql
`
	config := Config{}
	err := config.ParseConfig([]byte(data))

	assert.Error(t, err)
	assert.Empty(t, config.Services)
}

func TestNewServiceConfig(t *testing.T) {
	config := Config{}
	settings := map[string]string{
		"timeout": "10",
		"type": "mysql",
		"extra": "value",
	}
	serviceConfig, err := NewServiceConfig("name", settings, config)

	assert.Nil(t, err)
	assert.Equal(t, "name", serviceConfig.Name)
	assert.Equal(t, "mysql", serviceConfig.Type)
	assert.Equal(t, 10 * time.Second, serviceConfig.Timeout)
	assert.NotContains(t, serviceConfig.Settings, "timeout")
	assert.NotContains(t, serviceConfig.Settings, "type")
	assert.Contains(t, serviceConfig.Settings, "extra")
}