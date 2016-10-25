package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

func ShouldBeRun(a bot.Adapter, p *plugin.Plugin, m *bot.Message) bool {
	fmt.Printf("%+v", p)
	fmt.Printf("%+v", m)
	if p.RunOnlyOnChannels {
		return m.IsChannel
	}
	if p.RunOnlyOnDirectMessages {
		return m.IsDirectMessage
	}
	if p.RunOnlyOnMentions {
		return strings.Contains(m.Body, a.GetID())
	}
	return true
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
		// TODO: this is a little bit ugly
		plugin.RunOnlyOnChannels = pluginConfig.OnlyChannels
		plugin.RunOnlyOnDirectMessages = pluginConfig.OnlyDirectMessages
		plugin.RunOnlyOnMentions = pluginConfig.OnlyMentions

		if err != nil {
			log.Warningf("Error loading plugin (image: %s): %v", pluginConfig.Image, err)
			continue
		}
		defer plugin.Stop()
		plugins = append(plugins, plugin)
	}

	// TODO: this could be probably abstracted to be used as another adaptor (see Slack below)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debugf("HTTP POST received: %s", body)

		for _, p := range plugins {
			pluginResponse, err := p.Run(string(body))
			if err != nil {
				log.Error(err)
				continue
			}
			w.Write([]byte(pluginResponse + "\n"))
		}
	})

	host := fmt.Sprintf(":%d", config.Adapters.HTTP.Port)
	go func() { log.Fatal(http.ListenAndServe(host, nil)) }()

	log.Infof("HTTP adapter ready. Waiting for your POSTs at %s...", host)
	// --- END OF TODO ---

	// TODO: we shouldn't count that Slack is always configured
	slack, err := bot.New("slack", config.Adapters.Slack.Key)
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	log.Info("Slack adapter ready. Waiting for messages...")

	messagesCh, errorsCh := slack.Attach()
	for {
		select {
		case m := <-messagesCh:
			log.Debugf("Slack message received: %v", m)
			for _, p := range plugins {
				if !ShouldBeRun(slack, p, &m) {
					continue
				}
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
