package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Adapters struct {
		Slack struct {
			Key      string
			Channels []string
		}
	}
	Plugins []Plugin
}

type Plugin struct {
	Image       string
	Environment map[string]string
}

func NewFromFile(filePath string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err

	}
	return c, err
}
