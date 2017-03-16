package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/agonzalezro/botella/adapter"
	"github.com/agonzalezro/botella/config"
	"github.com/agonzalezro/botella/plugin"
)

func init() {
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}
}

func inferConfigPath() (string, error) {
	paths := []string{"botella.yml", "botella.yaml"}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("No %s file found!\n", strings.Join(paths, " or "))
}

// ensureVolumeHasMountPoint will return a list of volumes in the form `hostPath:containerPath` adding the `containerPath` in case it's missing.
func ensureVolumeHasMountPoint(vs []string) []string {
	volumes := []string{}
	for _, v := range vs {
		fragments := strings.Split(v, ":")
		if len(fragments) < 2 {
			v = fmt.Sprintf("%s:%s", fragments[0], fragments[0])
		}
		volumes = append(volumes, v)
	}
	return volumes
}

func loadPlugins(config *config.Config) ([]*plugin.Plugin, error) {
	var plugins []*plugin.Plugin
	for _, pluginConfig := range config.Plugins {
		plugin, err := plugin.New(
			pluginConfig.Image,
			pluginConfig.Environment,
			ensureVolumeHasMountPoint(pluginConfig.Volumes),
		)
		if err != nil {
			return nil, fmt.Errorf("Error loading plugin (image: %s): %v", pluginConfig.Image, err)
		}

		// TODO: this is a little bit ugly
		plugin.RunOnlyOnChannels = pluginConfig.OnlyChannels
		plugin.RunOnlyOnDirectMessages = pluginConfig.OnlyDirectMessages
		plugin.RunOnlyOnMentions = pluginConfig.OnlyMentions

		log.Infof("Plugin (%s) loaded.", pluginConfig.Image)
		log.Debugf("Plugin (%s) config: %+v", pluginConfig.Image, pluginConfig)
		plugins = append(plugins, plugin)
	}
	return plugins, nil
}

func loadAdapters(config *config.Config) ([]adapter.Adapter, error) {
	var adapters []adapter.Adapter
	for _, adapterConfig := range config.Adapters {
		adapter, err := adapter.New(adapterConfig.Name, adapterConfig.Environment)
		if err != nil {
			return nil, fmt.Errorf("Error loading adapter (%s): %v", adapterConfig.Name, err)
		}

		log.Infof("Adapter (%s) loaded.", adapterConfig.Name)
		log.Debugf("Adapter (%s) config: %+v", adapterConfig.Name, adapterConfig)
		adapters = append(adapters, adapter)
	}
	return adapters, nil
}

func listenAndReply(adapters []adapter.Adapter, plugins []*plugin.Plugin) {
	var wg sync.WaitGroup
	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh, os.Interrupt)

	for _, a := range adapters {
		wg.Add(1)

		stdinCh, stdoutCh, stderrCh := a.RunAndAttach()
		go func(a adapter.Adapter, stdinCh, stdoutCh chan adapter.Message, stderrCh chan error) {
			for {
				select {
				case m := <-stdinCh:
					log.Debugf("Message received: %+v", m)
					for _, p := range plugins {
						if !a.ShouldRun(p, &m) {
							log.Debugf("Not running plugin (%s) for: %+v", p.Image, m)
							continue
						}
						log.Debugf("Running plugin (%s) for: %+v", p.Image, m)

						stdout, stderr, err := p.Run(plugin.NewInput(m.Emitter, m.Receiver, m.Body))
						if err != nil {
							stderrCh <- err
							continue
						}
						stdout = strings.TrimSuffix(stdout, "\n")

						log.Debugf("Plugin (%s) response: %s", p.Image, stdout)
						if stderr != "" {
							log.Errorf("Plugin (%s) threw an error: %s", p.Image, stderr)
						}
						stdoutCh <- adapter.Message{Receiver: m.Receiver, Body: stdout}
					}
				case err := <-stderrCh:
					log.Error(err)
				case <-signalsCh:
					for i := 0; i < len(adapters); i++ {
						wg.Done()
					}
				}
			}
		}(a, stdinCh, stdoutCh, stderrCh)
	}

	wg.Wait()

	log.Info("Teardown...")
	for _, plugin := range plugins {
		plugin.Stop()
	}
}

func main() {
	configPath := flag.String("f", "", "Use a different file for the config. By default: botella.y{,a}ml")
	flag.Parse()

	if *configPath == "" {
		inferedPath, err := inferConfigPath()
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		configPath = &inferedPath
	}

	config, err := config.NewFromFile(*configPath)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}

	adapters, err := loadAdapters(config)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}

	plugins, err := loadPlugins(config)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}

	listenAndReply(adapters, plugins)
}
