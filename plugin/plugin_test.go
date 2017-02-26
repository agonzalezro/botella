package plugin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentArray(t *testing.T) {
	assert := assert.New(t)

	input := map[string]string{
		"PATH":   "path",
		"SECRET": "some-secret-key",
	}

	output := environmentAsArrayOfString("", input)
	assert.Equal(2, len(output))
	assert.EqualValues([]string{"PATH=path", "SECRET=some-secret-key"}, output)
}

// Test that if we define an empty value it tries to fallback to the system env vars.
func TestEnvironmentValuesFromEnvironmentVariables(t *testing.T) {
	assert := assert.New(t)

	err := os.Setenv("TEST_PLUGIN_SECRET", "value")
	assert.NoError(err)

	input := map[string]string{
		"secret": "",
	}

	output := environmentAsArrayOfString("test/plugin", input)
	assert.Equal(1, len(output))
	assert.Equal("secret=value", output[0])
}
