package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAValueFromAnEnvVar(t *testing.T) {
	assert := assert.New(t)

	err := os.Setenv("A_B_C_D", "value")
	assert.NoError(err)

	v, err := GetFromEnvOrFromMap("a/b-c", map[string]string{}, "d")
	assert.NoError(err)
	assert.Equal("value", v)
}

func TestAValueFromTheMap(t *testing.T) {
	kvs := map[string]string{
		"key": "value",
	}

	v, err := GetFromEnvOrFromMap("", kvs, "key")
	assert.NoError(t, err)
	assert.Equal(t, "value", v)
}

func TestNotFound(t *testing.T) {
	_, err := GetFromEnvOrFromMap("", nil, "not-found")
	assert.Error(t, err)
}
