package bot

import (
	"testing"

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
