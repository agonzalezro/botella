package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const validYAML = ` 
adapters:
  - name: slack
    environment:
      KEY: xxx
  - name: http
    environment:
      PORT: 8080

plugins:
  - image: agonzalezro/botella-test
    environment:
      KEY: this-is-a-secret
    volumes:
      - /Users/alex:/alex
    only_mentions: true
    only_direct_messages: true
    only_channels: true
`

func setup(assert *assert.Assertions) *os.File {
	content := []byte(validYAML)
	tmpfile, err := ioutil.TempFile("", "validYAML")
	assert.NoError(err)

	_, err = tmpfile.Write(content)
	assert.NoError(err)

	assert.NoError(tmpfile.Close())
	return tmpfile
}

func TestNewFromFile(t *testing.T) {
	assert := assert.New(t)
	tmpfile := setup(assert)

	config, err := NewFromFile(tmpfile.Name())
	assert.NoError(err)

	slack := config.Adapters[0]
	assert.Equal("slack", slack.Name)
	assert.Equal("xxx", slack.Environment["KEY"])

	http := config.Adapters[1]
	assert.Equal("http", http.Name)
	assert.Equal("8080", http.Environment["PORT"])

	assert.Equal(len(config.Plugins), 1)
	plugin := config.Plugins[0]
	assert.Equal("agonzalezro/botella-test", plugin.Image)
	assert.Equal(len(plugin.Environment), 1)
	assert.Equal("this-is-a-secret", plugin.Environment["KEY"])

	assert.Equal(len(plugin.Volumes), 1)
	assert.Equal("/Users/alex:/alex", plugin.Volumes[0])

	assert.Equal(true, plugin.OnlyDirectMessages)
	assert.Equal(true, plugin.OnlyMentions)
	assert.Equal(true, plugin.OnlyChannels)
}
