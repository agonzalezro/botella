package adapter

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"

	"github.com/agonzalezro/botella/plugin"
	"github.com/certifi/gocertifi"
)

const (
	rtmURLformatter = "https://slack.com/api/rtm.start?token=%s"
	wsURL           = "https://api.slack.com/"
)

type SlackAdapter struct {
	ws *websocket.Conn

	botID string
}

type SlackMessage struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (sm SlackMessage) isChannel() bool {
	// IDs in Slack seem to start with C for a Channel or with
	// G for a Group (private channel or group of people)
	return strings.HasPrefix(sm.Channel, "C") || strings.HasPrefix(sm.Channel, "G")
}

func (sm SlackMessage) isDirectMessage() bool {
	// IDs in Slack start with D for direct messages
	return strings.HasPrefix(sm.Channel, "D")
}

// TODO: this requires refactoring, it's tooooo long
func NewSlack(key string) (*SlackAdapter, error) {
	url := fmt.Sprintf(rtmURLformatter, key)

	cert_pool, err := gocertifi.CACerts()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: cert_pool},
	}
	client := http.Client{Transport: transport}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Received %d while connecting to Slack (expected 200)\n", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type Payload struct {
		Ok    bool
		Error string
		URL   string
		Self  struct {
			ID string
		}
	}
	var p Payload
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}
	if !p.Ok {
		return nil, errors.New(p.Error)
	}

	c, err := websocket.NewConfig(p.URL, wsURL)
	if err != nil {
		return nil, err
	}
	c.TlsConfig = &tls.Config{RootCAs: cert_pool}
	ws, err := websocket.DialConfig(c)
	if err != nil {
		return nil, err
	}

	return &SlackAdapter{ws: ws, botID: p.Self.ID}, nil
}

func (sa *SlackAdapter) ShouldRun(p *plugin.Plugin, m *Message) bool {
	if p.RunOnlyOnChannels {
		return m.IsChannel
	}
	if p.RunOnlyOnDirectMessages {
		return m.IsDirectMessage
	}
	if p.RunOnlyOnMentions {
		return strings.Contains(m.Body, sa.botID)
	}
	return true
}

func (sa *SlackAdapter) getSlackMessage() (*SlackMessage, error) {
	m := SlackMessage{}
	err := websocket.JSON.Receive(sa.ws, &m)
	return &m, err
}

func (sa *SlackAdapter) RunAndAttach() (chan Message, chan Message, chan error) {
	stdinCh := make(chan Message, 1)
	stdoutCh := make(chan Message, 1)
	stderrCh := make(chan error, 1)

	go func() {
		for {
			m, err := sa.getSlackMessage()
			if err != nil {
				stderrCh <- err
				continue
			}
			if m.Type == "message" {
				stdinCh <- Message{
					Emitter:         m.User,
					Receiver:        m.Channel,
					Body:            m.Text,
					IsChannel:       m.isChannel(),
					IsDirectMessage: m.isDirectMessage(),
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case m := <-stdoutCh:
				sm := SlackMessage{
					Type:    "message",
					Channel: m.Receiver,
					Text:    m.Body,
				}
				if err := websocket.JSON.Send(sa.ws, sm); err != nil {
					stderrCh <- err
				}
			}
		}
	}()

	return stdinCh, stdoutCh, stderrCh
}
