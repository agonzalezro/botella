package adapter

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/twinj/uuid"

	"github.com/agonzalezro/botella/plugin"
)

type HTTPAdapter struct {
	port int
}

func NewHTTP(port int) (*HTTPAdapter, error) {
	return &HTTPAdapter{port}, nil
}

func (HTTPAdapter) ShouldRun(_ *plugin.Plugin, _ *Message) bool {
	// This adapter doesn't have permissions
	return true
}

func (ha HTTPAdapter) RunAndAttach() (chan Message, chan Message, chan error) {
	stdinCh := make(chan Message, 1)
	stdoutCh := make(chan Message, 1)
	stderrCh := make(chan error, 1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			stderrCh <- err
			return
		}
		receiverID := uuid.NewV4().String()
		stdinCh <- Message{Receiver: receiverID, Body: string(body)}

		// FIXME: this loop isn't probably the best solution.
		// FIXME #2: it will just return one plugin response.
		for {
			m := <-stdoutCh
			if m.Receiver == receiverID {
				w.Write([]byte(m.Body + "\n"))
				return
			}
			stdoutCh <- m
		}
	})

	go func() {
		host := fmt.Sprintf(":%d", ha.port)
		stderrCh <- http.ListenAndServe(host, nil)
	}()

	return stdinCh, stdoutCh, stderrCh
}
