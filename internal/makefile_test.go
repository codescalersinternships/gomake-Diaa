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

func TestParseMakefile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		fileContent       string
		expectedError     error
		expectedAdjList   graph
		expectedTargsCmds commandMap
		failureMessage    string
	}{
		{
			name: "Valid format",
			fileContent: `run: build
	echo run
# comment
build:
	echo build`,
			expectedError: nil,
			expectedAdjList: graph{
				"run":   []string{"build"},
				"build": []string{},
			},
			expectedTargsCmds: commandMap{
				"run":   []string{"echo run"},
				"build": []string{"echo build"},
			},
			failureMessage: "got error while it's a valid format",
		},
		{
			name: "invalid format because of global command",
			fileContent: `
echo test
run:
	echo run`,
			expectedError:     errInvalidFormat,
			expectedAdjList:   nil,
			expectedTargsCmds: nil,
			failureMessage:    "failed to detect global command",
		},
		{
			name: "invalid format because of \t before target",
			fileContent: `
	run:
	echo run`,
			expectedError:     errInvalidFormat,
			expectedAdjList:   nil,
			expectedTargsCmds: nil,
			failureMessage:    "failed to detect \t before target",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.fileContent)
			gotAdjList, gotTargsCmds, err := parseMakefile(reader)

			assert.ErrorIs(t, err, tc.expectedError, tc.failureMessage)
			assert.Equal(t, tc.expectedAdjList, gotAdjList, tc.failureMessage)
			assert.Equal(t, tc.expectedTargsCmds, gotTargsCmds, tc.failureMessage)
		})
	}

}

func TestExtractTargetAndDeps(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		inputLine       string
		expTarget       string
		expDependencies []string
		failureMessage  string
	}{
		{
			name:            "Input not target",
			inputLine:       "echo run",
			expTarget:       "",
			expDependencies: nil,
			failureMessage:  "got target while input line is not target",
		},
		{
			name:            "Input is a target",
			inputLine:       "run: build",
			expTarget:       "run",
			expDependencies: []string{"build"},
			failureMessage:  "got not target while it's target",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotTarget, gotDeps := extractTargetAndDeps(tc.inputLine)

			assert.Equal(t, tc.expTarget, gotTarget, tc.failureMessage)
			assert.ElementsMatch(t, tc.expDependencies, gotDeps, tc.failureMessage)
		})
	}

}
