package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Adapters []Adapter
	Plugins  []Plugin
}

type Adapter struct {
	Name        string
	Environment map[string]string
}

type Plugin struct {
	Image              string
	Environment        map[string]string
	OnlyChannels       bool `yaml:"only_channels"`
	OnlyDirectMessages bool `yaml:"only_direct_messages"`
	OnlyMentions       bool `yaml:"only_mentions"`
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

func (c Config) EnvironmentFor(adapterName string) (map[string]string, error) {
	for _, adapterConfig := range c.Adapters {
		if adapterConfig.Name == adapterName {
			return adapterConfig.Environment, nil
		}
	}
	return nil, fmt.Errorf("Configuration for adapter %s not found", adapterName)
}
