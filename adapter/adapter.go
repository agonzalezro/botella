package adapter

import (
	"fmt"
	"strconv"

	"github.com/agonzalezro/ava/plugin"
	"github.com/agonzalezro/ava/utils"
)

type Message struct {
	// Emitter is the emitter of the message, in the case of Slack is clear, in
	// some other cases as for example the http adapter, it doesn't need to be set
	Emitter string
	// Receiver is the receiver of the message, for example: a channel ID
	Receiver string
	Body     string

	IsChannel       bool
	IsDirectMessage bool
}

type Adapter interface {
	RunAndAttach() (stdin chan Message, stdout chan Message, stderr chan error)
	ShouldRun(*plugin.Plugin, *Message) bool
}

func New(adapterName string, environment map[string]string) (Adapter, error) {
	switch adapterName {
	case "slack":
		key, err := utils.GetFromEnvOrFromMap(adapterName, environment, "key")
		if err != nil {
			return nil, err
		}
		return NewSlack(key)
	case "http":
		port, err := utils.GetFromEnvOrFromMap(adapterName, environment, "port")
		if err != nil {
			return nil, err
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
