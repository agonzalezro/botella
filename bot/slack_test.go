package bot

import (
	"fmt"
	"testing"

	"github.com/agonzalezro/ava/plugin"
	"github.com/stretchr/testify/assert"
)

func TestChannelOrDirectMessage(t *testing.T) {
	assert := assert.New(t)

	type result struct {
		IsChannel, IsDirectMessage bool
	}

	cases := map[string]result{
		"C1PP69WMA": result{true, false},
		"G2TGW8ETA": result{true, false},
		"D1PQQAGTZ": result{false, true},
		"AABBCCDD1": result{false, false},
	}

	for in, out := range cases {
		sm := SlackMessage{Channel: in}
		assert.Equal(out.IsChannel, sm.isChannel())
		assert.Equal(out.IsDirectMessage, sm.isDirectMessage())
	}
}

func TestIfPluginShouldBeRun(t *testing.T) {
	assert := assert.New(t)

	adapter := SlackAdapter{botID: "test-id"}

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
		m := Message{
			IsChannel:       in.isChannel,
			IsDirectMessage: in.isDirectMessage,
			Body:            in.body,
		}

		assert.Equal(
			expected,
			adapter.ShouldRun(&p, &m),
			fmt.Sprintf("%+v %+v %+v", adapter, p, m))
	}
}
