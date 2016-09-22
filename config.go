package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config holds configuration for Starter
type Config struct {
	APIURL        string
	template_path string
	use_registry  bool
}

// ReadFromFile reads config from a file
func ReadFromFile(configFile string) (*Config, error) {
	var config *Config
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	config.SetDefaults()

	return config, nil
}

// SetDefaults sets default values for unset config items
func (c *Config) SetDefaults() error {
	if c.APIURL == "" {
		c.APIURL = "0.0.0.0:9090"
	}
	c.use_registry = false
	return nil
}
