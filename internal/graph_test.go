package makefile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDependencyGraph(t *testing.T) {
	t.Parallel()
	dg := NewDependencyGraph()

	assert.NotNil(t, dg, "expected dependency graph to be initialized, but was nil")

	assert.NotNil(t, dg.adjacencyList, "expected adjacencyList to be initialized, but was nil")

	assert.NotNil(t, dg.targetToCommands, "expected TargetToCommands to be initialized, but was nil")
}

func TestSetAdjacencyList(t *testing.T) {
	t.Parallel()
	dg := NewDependencyGraph()

	adjList := graph{
		"run":   []string{"build"},
		"build": []string{},
	}

	dg.SetAdjacencyList(adjList)

	assert.Equal(t, adjList, dg.adjacencyList, "fail to set adjacency list")
}

func TestSetTargetToCommands(t *testing.T) {
	t.Parallel()

	dg := NewDependencyGraph()

	targetToCommands := commandMap{
		"run":     []string{"npm run start", "npm run start:dev"},
		"install": []string{"npm install"},
	}

	dg.SetTargetToCommands(targetToCommands)

	assert.Equal(t, targetToCommands, dg.targetToCommands, "fail to set target to commands")
}

func TestHasCircularDependency(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		adjList        graph
		expectedError  error
		failureMessage string
	}{
		{
			name: "Circular dependencies exist",
			adjList: graph{
				"run":   []string{"build"},
				"build": []string{"run"},
			}, expectedError: errCycleDetected,
			failureMessage: "fail to detect circular dependency in dependency graph with cycles",
		}, {
			name: "No circular dependencies",
			adjList: graph{
				"run":   []string{"build"},
				"build": []string{},
			},
			expectedError:  nil,
			failureMessage: "circular dependency detected in dependency graph with no cycles",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			dg := NewDependencyGraph()
			dg.SetAdjacencyList(tc.adjList)

			err := dg.checkCircularDependency()

			assert.ErrorIs(t, err, tc.expectedError, tc.failureMessage)

		})
	}

}

func TestCheckMissingDependencies(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		adjList        graph
		expected       []string
		failureMessage string
	}{
		{
			name: "Has no missing dependencies",
			adjList: graph{
				"run":   []string{"build"},
				"build": []string{},
			},
			expected:       []string{},
			failureMessage: "got missing dependencies while shouldn't",
		}, {
			name: "Has missing dependencies",
			adjList: graph{
				"run": []string{"build", "make"},
			},
			expected:       []string{"build", "make"},
			failureMessage: "got missing dependencies not as expected",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dg := NewDependencyGraph()

			dg.SetAdjacencyList(tc.adjList)

			got := dg.checkMissingDependencies()

			assert.ElementsMatch(t, tc.expected, got, tc.failureMessage)
		})
	}
}

func TestExecuteTargetKAndItsDeps(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		adjList        graph
		targetCommands commandMap
		failureMessage string
		expectedError  error
	}{
		{
			name:           "Target doesn't exist",
			adjList:        graph{},
			targetCommands: commandMap{},
			failureMessage: "fail to detect that target doesn't exist",
			expectedError:  errTargetDoesnotExist,
		}, {
			name: "Target exist, should exec commands",
			adjList: graph{
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
			dg := NewDependencyGraph()

			dg.SetAdjacencyList(tc.adjList)

			dg.SetTargetToCommands(tc.targetCommands)

			err := dg.ExecuteTargetKAndItsDeps("run")

			assert.ErrorIs(t, err, tc.expectedError, tc.failureMessage)
		})
	}

}

func TestExecuteTasksInDependencyOrder(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		adjList  graph
		commands commandMap
		target   string
		expected string
	}{
		{
			adjList: graph{
				"run":   []string{"build", "print"},
				"build": []string{"exec"},
				"exec":  []string{},
				"print": []string{"exec"},
			},
			commands: commandMap{
				"run":   []string{"echo run"},
				"build": []string{"echo build"},
				"exec":  []string{"echo exec"},
				"print": []string{"echo print"},
			},
			target: "run",
			expected: `echo exec
exec
echo build
build
echo print
print
echo run
run
`,
		},
	}

	for idx, tc := range testCases {
		t.Run("should execute commands in right order", func(t *testing.T) {

			dg := NewDependencyGraph()
			dg.SetAdjacencyList(tc.adjList)
			dg.SetTargetToCommands(tc.commands)

			visited := make(map[string]bool)
			got, err := dg.executeTasksInDependencyOrder(tc.target, visited)

			assert.Nil(t, err, "failed to execute commands while it shouldn't in test #%d", idx+1)
			assert.Equal(t, tc.expected, got, "failed to execute in the right dependencies order in test #%d", idx+1)
		})
	}
}
func TestExecuteCommandsForTargetK(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		adjList        graph
		targetCommands commandMap
		failureMessage string
		expectedError  error
		target         string
	}{
		{
			name: "Target has no commands",
			adjList: graph{
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
			adjList: graph{
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
			dg := NewDependencyGraph()

			dg.SetTargetToCommands(tc.targetCommands)

			_, err := dg.executeCommandsForTargetK(tc.target)

			assert.ErrorIs(t, err, tc.expectedError, tc.failureMessage)

		})
	}
}
