package makefile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDependencyGraph(t *testing.T) {
	t.Parallel()
	dg := NewDependencyGraph(graph{}, commandMap{})

	assert.NotNil(t, dg, "expected dependency graph to be initialized, but was nil")

	assert.NotNil(t, dg.adjacencyList, "expected adjacencyList to be initialized, but was nil")

	assert.NotNil(t, dg.targetToCommands, "expected TargetToCommands to be initialized, but was nil")
}

func TestHasCircularDependency(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		adjacencyList  graph
		expectedError  error
		failureMessage string
	}{
		{
			name: "Circular dependencies exist",
			adjacencyList: graph{
				"run":   []string{"build"},
				"build": []string{"run"},
			}, expectedError: errCycleDetected,
			failureMessage: "fail to detect circular dependency in dependency graph with cycles",
		}, {
			name: "No circular dependencies",
			adjacencyList: graph{
				"run":   []string{"build"},
				"build": []string{},
			},
			expectedError:  nil,
			failureMessage: "circular dependency detected in dependency graph with no cycles",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			dg := NewDependencyGraph(tc.adjacencyList, commandMap{})

			err := dg.checkCircularDependency()

			assert.ErrorIs(t, err, tc.expectedError, tc.failureMessage)

		})
	}

}

func TestCheckMissingDependencies(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		adjacencyList  graph
		expected       []string
		failureMessage string
	}{
		{
			name: "Has no missing dependencies",
			adjacencyList: graph{
				"run":   []string{"build"},
				"build": []string{},
			},
			expected:       []string{},
			failureMessage: "got missing dependencies while shouldn't",
		}, {
			name: "Has missing dependencies",
			adjacencyList: graph{
				"run": []string{"build", "make"},
			},
			expected:       []string{"build", "make"},
			failureMessage: "got missing dependencies not as expected",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			dg := NewDependencyGraph(tc.adjacencyList, commandMap{})

			got := dg.checkMissingDependencies()

			assert.ElementsMatch(t, tc.expected, got, tc.failureMessage)
		})
	}
}

func TestExecuteTargetAndItsDeps(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		adjacencyList  graph
		targetCommands commandMap
		failureMessage string
		expectedError  error
	}{
		{
			name:           "Target doesn't exist",
			adjacencyList:  graph{},
			targetCommands: commandMap{},
			failureMessage: "fail to detect that target doesn't exist",
			expectedError:  errTargetDoesnotExist,
		}, {
			name: "Target exist, should exec commands",
			adjacencyList: graph{
				"run": []string{},
			},
			targetCommands: commandMap{
				"run": []string{"echo test"},
			},
			failureMessage: "fail to execute commands while it shouldn't",
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			dg := NewDependencyGraph(tc.adjacencyList, tc.targetCommands)

			err := dg.ExecuteTargetAndItsDeps("run")

			assert.ErrorIs(t, err, tc.expectedError, tc.failureMessage)
		})
	}

}

func TestGetTasksOrder(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		adjacencyList   graph
		targetsCommands commandMap
		target          string
		expectedOrder   []string
	}{
		{
			adjacencyList: graph{
				"run":   []string{"build", "print"},
				"build": []string{"exec"},
				"exec":  []string{},
				"print": []string{"exec"},
			},
			targetsCommands: commandMap{
				"run":   []string{"echo run"},
				"build": []string{"echo build"},
				"exec":  []string{"echo exec"},
				"print": []string{"echo print"},
			},
			target:        "run",
			expectedOrder: []string{"exec", "build", "print", "run"},
		},
	}

	for idx, tc := range testCases {
		t.Run("should execute commands in right order", func(t *testing.T) {

			dg := NewDependencyGraph(tc.adjacencyList, tc.targetsCommands)

			visited := make(map[string]bool)
			gotOrder := dg.getTasksOrder(tc.target, visited)

			assert.Equal(t, tc.expectedOrder, gotOrder, "failed to execute in the right dependencies order in test #%d", idx+1)
		})
	}
}
func TestExecuteCommandsForTarget(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		adjacencyList  graph
		targetCommands commandMap
		failureMessage string
		expectedError  error
		target         string
	}{
		{
			name: "Target has no commands",
			adjacencyList: graph{
				"run": []string{},
			},
			targetCommands: commandMap{
				"run": []string{},
			},
			failureMessage: "failed to detect that there is no commands",
			expectedError:  errTargetHasNoCommands,
			target:         "run",
		}, {
			name: "Target has commands",
			adjacencyList: graph{
				"run": []string{},
			},
			targetCommands: commandMap{
				"run": []string{"echo test"},
			},
			failureMessage: "failed to detect the exists command",
			expectedError:  nil,
			target:         "run",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dg := NewDependencyGraph(graph{}, tc.targetCommands)

			err := dg.executeCommandsForTarget(tc.target)

			assert.ErrorIs(t, err, tc.expectedError, tc.failureMessage)

		})
	}
}
