package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentArray(t *testing.T) {
	assert := assert.New(t)

	input := map[string]string{
		"PATH":   "path",
		"SECRET": "some-secret-key",
	}
	output := environmentAsArrayOfString(input)

	_ = assert
	_ = output
	//assert.Equal("PATH=path", output[0])
	//assert.Equal("SECRET=some-secret-key", output[1])
}
