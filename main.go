package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/agonzalezro/ava/bot"
	"github.com/agonzalezro/ava/config"
	"github.com/agonzalezro/ava/plugin"
)

func init() {
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}
}

func inferConfigPath() (string, error) {
	paths := []string{"ava.yml", "ava.yaml"}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("No %s file found!", strings.Join(paths, " or "))
}

func main() {
	configPath, err := inferConfigPath()
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}

	config, err := config.NewFromFile(configPath)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}

	var plugins []*plugin.Plugin
	for _, pluginConfig := range config.Plugins {
		plugin, err := plugin.New(
			pluginConfig.Image,
			pluginConfig.Environment,
		)
		if err != nil {
			log.Warningf("Error loading plugin (image: %s): %v", pluginConfig.Image, err)
			continue
		}
		defer plugin.Stop()
		plugins = append(plugins, plugin)
	}

	// TODO: in the future there should be more adaptors
	slack, err := bot.New("slack", config.Adapters.Slack.Key)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	log.Info("Slack adaptor ready. Waiting for messages...")

	messagesCh, errorsCh := slack.Attach()
	for {
		select {
		case m := <-messagesCh:
			log.Debugf("Slack message received: %v", m)
			for _, p := range plugins {
				pluginResponse, err := p.Run(m.Body)
				if err != nil {
					errorsCh <- err
					continue
				}
				pluginResponse = strings.TrimSuffix(pluginResponse, "\n")

				log.Debugf("Plugin (%s) response: %s", p.Image, pluginResponse)
				slack.Send(bot.Message{Channel: m.Channel, Body: pluginResponse})
			}
		case err := <-errorsCh:
			log.Error(err)
		}
	}
}
