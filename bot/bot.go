package bot

import "fmt"

type Message struct {
	Channel string
	Body    string
}

type Adaptor interface {
	Attach() (chan Message, chan error)
	Send(Message) error
}

func New(adaptor, key string) (Adaptor, error) {
	switch adaptor {
	case "slack":
		return NewSlack(key)
	default:
		return nil, fmt.Errorf("Adaptor '%s' not found", adaptor)
	}
}
