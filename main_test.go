package main

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
)

func TestInitConfig(t *testing.T) {
	initConfig()

	assert.Equal(t, viper.GetInt("port"), 8080)
}
