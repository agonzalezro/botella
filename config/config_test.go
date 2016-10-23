package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const validYAML = ` 
adapters:
  slack:
    key: xxx
    channels:
      - general
    name: Ava González
    pic: url

plugins:
  - image: agonzalezro/ava-test
    environment:
      KEY: this-is-a-secret
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

	slack := config.Adapters.Slack
	assert.Equal("xxx", slack.Key)
	assert.Equal([]string{"general"}, slack.Channels)

	assert.Equal(len(config.Plugins), 1)
	plugin := config.Plugins[0]
	assert.Equal("agonzalezro/ava-test", plugin.Image)
	assert.Equal(len(plugin.Environment), 1)
	assert.Equal("this-is-a-secret", plugin.Environment["KEY"])
}
