package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/agonzalezro/ava/plugin"

	"golang.org/x/net/websocket"
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

func NewSlack(key string) (*SlackAdapter, error) {
	url := fmt.Sprintf(rtmURLformatter, key)
	resp, err := http.Get(url)
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

	ws, err := websocket.Dial(p.URL, "", wsURL)
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
					Channel:         m.Channel,
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
					Channel: m.Channel,
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
