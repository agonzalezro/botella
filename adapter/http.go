package adapter

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/agonzalezro/ava/plugin"
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
		stdinCh <- Message{Body: string(body)}

		// FIXME: this will just return one plugin response!
		// And actually it doesn't even need to be the one that ask for it :D
		m := <-stdoutCh
		w.Write([]byte(m.Body + "\n"))
	})

	go func() {
		host := fmt.Sprintf(":%d", ha.port)
		stderrCh <- http.ListenAndServe(host, nil)
	}()

	return stdinCh, stdoutCh, stderrCh
}
