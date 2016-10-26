package bot

import "fmt"

type Message struct {
	Channel         string
	Body            string
	IsChannel       bool
	IsDirectMessage bool
}

type Adapter interface {
	GetID() string

	Attach() (chan Message, chan error)
	Send(Message) error
}

func New(adapter, key string) (Adapter, error) {
	switch adapter {
	case "slack":
		return NewSlack(key)
	default:
		return nil, fmt.Errorf("Adapter '%s' not found\n", adapter)
	}
}
