package makefile

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMakefile(t *testing.T) {

	t.Parallel()

	testCases := []struct {
		name              string
		filePath          string
		pathValidation    string
		failureMessage    string
		needToWriteFile   bool
		shouldReturnError bool
	}{
		{
			name:              "Invalid file path",
			filePath:          "file",
			failureMessage:    "should return error no such file or directory but got a file",
			shouldReturnError: true,
			needToWriteFile:   false,
		}, {
			name:              "Valid file path",
			filePath:          "Makefile",
			failureMessage:    "should return nil but got error",
			shouldReturnError: false,
			needToWriteFile:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if tc.needToWriteFile {
				os.WriteFile(path.Join(dir, tc.filePath), []byte{}, 0644)
			}

			_, _, err := ReadMakefile(path.Join(dir, tc.filePath))
			if tc.shouldReturnError {
				assert.NotNil(t, err, tc.failureMessage)
			} else {
				assert.Nil(t, err, tc.failureMessage)
			}
		})
	}
}

func TestParseMakefile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                    string
		fileContent             string
		expectedError           error
		expectedAdjacencyList   graph
		expectedTargetsCommands commandMap
		failureMessage          string
	}{
		{
			name: "Valid format",
			fileContent: `run: build
	echo run
# comment
build:
	echo build`,
			expectedError: nil,
			expectedAdjacencyList: graph{
				"run":   []string{"build"},
				"build": []string{},
			},
			expectedTargetsCommands: commandMap{
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
			expectedError:           errInvalidFormat,
			expectedAdjacencyList:   nil,
			expectedTargetsCommands: nil,
			failureMessage:          "failed to detect global command",
		},
		{
			name: "invalid format because of \t before target",
			fileContent: `
	run:
	echo run`,
			expectedError:           errInvalidFormat,
			expectedAdjacencyList:   nil,
			expectedTargetsCommands: nil,
			failureMessage:          "failed to detect \t before target",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.fileContent)
			gotAdjacencyList, gotTargetsToCommands, err := parseMakefile(reader)

			fmt.Println(gotTargetsToCommands)
			assert.ErrorIs(t, err, tc.expectedError, tc.failureMessage)
			assert.Equal(t, tc.expectedAdjacencyList, gotAdjacencyList, tc.failureMessage)
			assert.Equal(t, tc.expectedTargetsCommands, gotTargetsToCommands, tc.failureMessage)
		})
	}

}

func TestIsTarget(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		inputLine      string
		expected       bool
		failureMessage string
	}{
		{
			name:           "Input not target",
			inputLine:      "echo run",
			failureMessage: "got target while input line is not target",
			expected:       false,
		},
		{
			name:           "Input is a target",
			inputLine:      "run: build",
			failureMessage: "got not target while it's target",
			expected:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isTarget(tc.inputLine)

			assert.Equal(t, tc.expected, got, tc.failureMessage)
		})
	}

}

func TestextractTargetAndDeps(t *testing.T) {

	testCases := []struct {
		name            string
		inputLine       string
		expTarget       string
		expDependencies []string
	}{
		{
			name:            "test extract target and dependencies",
			inputLine:       "run: build",
			expTarget:       "run",
			expDependencies: []string{"build"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotTarget, gotDependencies := extractTargetAndDeps(tc.inputLine)

			assert.ElementsMatch(t, tc.expDependencies, gotDependencies)

			assert.Equal(t, tc.expTarget, gotTarget)
		})
	}

}
