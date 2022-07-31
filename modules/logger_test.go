package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Log(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		expected []string
	}{
		{
			name:     "with no message",
			msg:      "",
			expected: []string{"dogs"},
		},
		{
			name:     "with a message",
			msg:      "cats",
			expected: []string{"cats", "dogs"},
		},
		{
			name:     "with same message",
			msg:      "dogs",
			expected: []string{"dogs", "dogs"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger("logger")
			logger.Log("dogs")

			logger.Log(tt.msg)
			actual := logger.Messages

			assert.IsType(t, tt.expected, actual)
		})
	}
}
