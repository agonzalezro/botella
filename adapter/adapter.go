package adapter

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/agonzalezro/ava/plugin"
)

type Message struct {
	Channel         string
	Body            string
	IsChannel       bool
	IsDirectMessage bool
}

type Adapter interface {
	GetID() string
	RunAndAttach() (stdin chan Message, stdout chan Message, stderr chan error)
	ShouldRun(*plugin.Plugin, *Message) bool
}

func New(adapterName string, environment map[string]string) (Adapter, error) {
	switch adapterName {
	case "slack":
		key, ok := environment["key"]
		if !ok {
			return nil, errors.New("key field is mandatory in environment conf for Slack")
		}
		return NewSlack(key)
	case "http":
		port, ok := environment["port"]
		if !ok {
			return nil, errors.New("port is mandatory in environment conf for HTTP")
		}
		iport, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("port in HTTP adapter should be an integer, it's: %s", port)
		}
		return NewHTTP(iport)
	default:
		return nil, fmt.Errorf("Adapter '%s' not found\n", adapterName)
	}
}
