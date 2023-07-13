package makefile

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMakefile(t *testing.T) {

	t.Parallel()

	testCases := []struct {
		name           string
		filePath       string
		pathValidation string
		failureMessage string
	}{
		{
			name:           "Invalid file path",
			filePath:       "file",
			pathValidation: "invalid",
			failureMessage: "should return error no such file or directory but got a file",
		}, {
			name:           "Valid file path",
			filePath:       "Makefile",
			pathValidation: "",
			failureMessage: "should return nil but got error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			os.WriteFile(path.Join(dir, tc.filePath), []byte{}, 0644)

			_, _, err := ReadMakefile(path.Join(dir, tc.filePath, tc.pathValidation))
			if tc.pathValidation == "" {
				assert.Nil(t, err, tc.failureMessage)
			} else {
				assert.NotNil(t, err, tc.failureMessage)
			}
		})
	}
}
