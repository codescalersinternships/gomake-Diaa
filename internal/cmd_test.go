package makefile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecCommand(t *testing.T) {

	testCases := []struct {
		name         string
		command      string
		expectingErr bool
	}{
		{
			name:         "invalid binary",
			command:      "binary test",
			expectingErr: true,
		}, {
			name:         "valid binary",
			command:      "echo test",
			expectingErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			err := execCommand(tc.command)

			if tc.expectingErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

		})
	}
}
