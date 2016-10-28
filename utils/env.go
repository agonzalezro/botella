package utils

import (
	"fmt"
	"os"
	"strings"
)

func sanitizePrefix(prefix string) string {
	replacements := []string{"/", "-"}

	for _, replacement := range replacements {
		prefix = strings.Replace(prefix, replacement, "_", -1)
	}
	return prefix
}

// GetFromEnvOrFromMap will try to find the the key in the provided map
// or in an environment variable.
//
// If the environment variable is set, it has preference. The environment var
// will be queried all uppercased.
//
// Note: the / char in the prefix will be transformed to _
func GetFromEnvOrFromMap(prefix string, kvs map[string]string, k string) (string, error) {
	envVar := strings.ToUpper(fmt.Sprintf("%s_%s", sanitizePrefix(prefix), k))
	v := os.Getenv(envVar)
	if v != "" {
		return v, nil
	}

	if v, ok := kvs[k]; ok {
		return v, nil
	}

	return "", fmt.Errorf("'%s' env var was not found, neither the key '%s' in the configuration.", envVar, k)
}
