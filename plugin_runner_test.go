package main


import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewPluginRunner(t *testing.T) {
	pr, err := NewPluginRunner("location")

	assert.Nil(t, err)
	assert.Equal(t, pr.location, "location")
}
