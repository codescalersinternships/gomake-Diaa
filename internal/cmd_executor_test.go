package makefile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecCommand(t *testing.T) {

	testCases := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "valid command with @",
			command:  "@echo test",
			expected: "test\n",
		},
		{
			name:     "valid command without @",
			command:  "echo test",
			expected: "echo test\ntest\n",
		},
		{
			name:    "invalid binary",
			command: "binary test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			got, err := execCommand(tc.command)

			if err != nil {
				assert.NotNil(t, err, "expecting binary error but got no error")
			} else {
				assert.Equal(t, tc.expected, got, "output doesn't match")
			}
		})
	}
}
