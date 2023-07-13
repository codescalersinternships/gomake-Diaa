package makefile

import (
	"os"
	"path"
	"strings"
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

// func TestParseMakefile(t *testing.T) {
// 	t.Parallel()

// 	testCases := []struct {
// 		name           string
// 		fileContent    string
// 		expectedError  error
// 		expectedAdjList Graph
// 		expectedTargsCmds CommandMap
// 		failureMessage string
// 	}{
// 		{
// 			name: "Valid format",
// 			fileContent: `run: build
// 	echo run
// # comment
// build:
// 	echo build`,
// 			expectedError:  nil,
// 			expectedAdjList: Graph{
// 				"run":[]string{"build"},
// 				"build":[] string{},
// 			},
// 			expectedTargsCmds: CommandMap{
// 				"run":"echo run",
// 				"build":"echo build",
// 			},
// 			failureMessage: "got error while it's a valid format",
// 		},
// 		{
// 			name: "invalid format because of global command",
// 			fileContent: `
// echo test
// run:
// 	echo run`,
// 			expectedError:  ErrInvalidFormat,
// 			expectedAdjList: nil,
// 			expectedTargsCmds: nil,
// 			failureMessage: "failed to detect global command",
// 		},
// 		{
// 			name: "invalid format because of \t before target",
// 			fileContent: `
// 	run:
// 	echo run`,
// 			expectedError:  ErrInvalidFormat,
// 			expectedAdjList: nil,
// 			expectedTargsCmds: nil,
// 			failureMessage: "failed to detect \t before target",
// 		},
// 	}

// 	for _,tc := range testCases{
// 		t.Run(tc.name,func(t *testing.T) {
// 			reader := strings.NewReader(tc.fileContent)
// 			gotAdjList,gotTargsCmds,err:= ParseMakefile(reader)

// 			assert.
// 		})
// 	}

// }
