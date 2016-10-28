package main

import (
	"testing"

	"github.com/agonzalezro/ava/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadPlugins(t *testing.T) {
	assert := assert.New(t)

	image := "busybox"
	pluginConfig := config.Plugin{
		Image:              image, // TODO: not mocked, we will need Docker running
		OnlyChannels:       true,
		OnlyDirectMessages: true,
		OnlyMentions:       true,
	}
	config := config.Config{
		Plugins: []config.Plugin{pluginConfig},
	}

	plugins, err := loadPlugins(&config)
	assert.NoError(err)
	assert.Equal(1, len(plugins))

	p := plugins[0]
	defer p.Stop() // TODO: ugly but better than nothing

	assert.Equal(image, p.Image)
	assert.True(p.RunOnlyOnChannels)
	assert.True(p.RunOnlyOnDirectMessages)
	assert.True(p.RunOnlyOnMentions)
}

func TestLoadPluginThatErrors(t *testing.T) {
	pluginConfig := config.Plugin{Image: "this-plugin-does-not-exist"}
	config := config.Config{
		Plugins: []config.Plugin{pluginConfig},
	}

	_, err := loadPlugins(&config)
	assert.Error(t, err)
}

func TestLoadAdapters(t *testing.T) {
	assert := assert.New(t)

	adapterConfig := config.Adapter{
		Name: "http",
		Environment: map[string]string{
			"port": "56789",
		},
	}
	config := config.Config{
		Adapters: []config.Adapter{adapterConfig},
	}

	adapters, err := loadAdapters(&config)
	assert.NoError(err)
	assert.Equal(1, len(adapters))
}

func TestLoadAdaptersNotFound(t *testing.T) {
	adapterConfig := config.Adapter{Name: "does-not-exist"}
	config := config.Config{
		Adapters: []config.Adapter{adapterConfig},
	}

	_, err := loadAdapters(&config)
	assert.Error(t, err)
}
