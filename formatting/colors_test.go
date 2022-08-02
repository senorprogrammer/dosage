package formatting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ColorForState(t *testing.T) {
	tests := []struct {
		name     string
		state    string
		body     string
		expected string
	}{
		{
			name:     "with active",
			state:    "active",
			body:     "cats",
			expected: "[green:]cats[-:]",
		},
		{
			name:     "with off",
			state:    "off",
			body:     "cats",
			expected: "[red:]cats[-:]",
		},
		{
			name:     "with offline",
			state:    "offline",
			body:     "cats",
			expected: "[red:]cats[-:]",
		},
		{
			name:     "with online",
			state:    "online",
			body:     "cats",
			expected: "[green:]cats[-:]",
		},
		{
			name:     "with unknown state",
			state:    "dogs",
			body:     "cats",
			expected: "[white:]cats[-:]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func Test_Bold(t *testing.T) {
	assert.Equal(t, "[::b]cats[::-]", Bold("cats"))
}
