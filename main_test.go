package main

import (
	"fmt"
	"testing"

	"github.com/agonzalezro/ava/bot"
	"github.com/agonzalezro/ava/plugin"
	"github.com/stretchr/testify/assert"
)

type testAdaptor struct{}

func (testAdaptor) GetID() string                          { return "test-id" }
func (testAdaptor) Attach() (chan bot.Message, chan error) { return nil, nil }
func (testAdaptor) Send(bot.Message) error                 { return nil }

func TestIfPluginShouldBeRun(t *testing.T) {
	assert := assert.New(t)

	adaptor := testAdaptor{}

	type c struct {
		runOnlyOnChannels, runOnlyOnDirectMessages, runOnlyOnMentions bool
		isChannel, isDirectMessage                                    bool
		body                                                          string
	}
	cases := map[c]bool{
		c{runOnlyOnChannels: true, isChannel: true}:              true,
		c{runOnlyOnChannels: true, isChannel: false}:             false,
		c{runOnlyOnDirectMessages: true, isDirectMessage: true}:  true,
		c{runOnlyOnDirectMessages: true, isDirectMessage: false}: false,
		c{runOnlyOnMentions: true, body: "<test-id> run this"}:   true,
		c{runOnlyOnMentions: true, body: "not mentioned"}:        false,
		c{}: true,
	}

	for in, expected := range cases {
		p := plugin.Plugin{
			RunOnlyOnChannels:       in.runOnlyOnChannels,
			RunOnlyOnDirectMessages: in.runOnlyOnDirectMessages,
			RunOnlyOnMentions:       in.runOnlyOnMentions,
		}
		m := bot.Message{
			IsChannel:       in.isChannel,
			IsDirectMessage: in.isDirectMessage,
			Body:            in.body,
		}

		assert.Equal(
			expected,
			ShouldBeRun(&adaptor, p, m),
			fmt.Sprintf("%+v %+v %+v", adaptor, p, m))
	}

	// if run only on DMs but it's a channel
	// if run only on mentioned but it isn't mentioned
	// if no run only at all
}
