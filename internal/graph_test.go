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

	adjList := make(Graph)

	adjList["run"] = []string{"build"}
	adjList["build"] = []string{}

	dg.SetAdjacencyList(adjList)

	assert.Equal(t, adjList, dg.adjacencyList, "fail to set adjacency list")
}

func TestSetTargetToCommands(t *testing.T) {
	t.Parallel()

	dg := NewDependencyGraph()
	targetToCommands := make(CommandMap)
	targetToCommands["run"] = []string{"npm run start", "npm run start:dev"}
	targetToCommands["install"] = []string{"npm install"}

	dg.SetTargetToCommands(targetToCommands)

	assert.Equal(t, targetToCommands, dg.targetToCommands, "fail to set target to commands")
}

func TestHasCircularDependency(t *testing.T) {
	t.Parallel()

	t.Run("Circular dependencies exist", func(t *testing.T) {
		dg := NewDependencyGraph()

		adjList := make(Graph)

		adjList["run"] = []string{"build"}
		adjList["build"] = []string{"run"}

		dg.SetAdjacencyList(adjList)

		err := dg.CheckCircularDependency()

		assert.ErrorIs(t, err, ErrCycleDetected, "fail to detect circular dependency in dependency graph with cycles")

	})

	t.Run("No circular dependencies", func(t *testing.T) {
		dg := NewDependencyGraph()

		adjList := make(Graph)

		adjList["run"] = []string{"build"}
		adjList["build"] = []string{}

		dg.SetAdjacencyList(adjList)

		err := dg.CheckCircularDependency()

		assert.Nil(t, err, "circular dependency detected in dependency graph with no cycles")
	})
}

func TestCheckMissingDependencies(t *testing.T) {
	t.Parallel()

	t.Run("Has no missing dependencies", func(t *testing.T) {
		dg := NewDependencyGraph()

		adjList := make(Graph)

		adjList["run"] = []string{"build"}
		adjList["build"] = []string{}

		dg.SetAdjacencyList(adjList)

		got := dg.CheckMissingDependencies()

		want := []string{}

		assert.Equal(t, want, got, "got missing dependencies while shouldn't")
	})

	t.Run("Has missing dependencies", func(t *testing.T) {
		dg := NewDependencyGraph()

		adjList := make(Graph)

		adjList["run"] = []string{"build", "make"}

		dg.SetAdjacencyList(adjList)

		got := dg.CheckMissingDependencies()

		want := []string{"build", "make"}

		assert.ElementsMatch(t, want, got, "got missing dependencies not as expected")
	})
}

func TestExecuteTargetKAndItsDeps(t *testing.T) {
	t.Parallel()

	t.Run("Target doesn't exist", func(t *testing.T) {
		dg := NewDependencyGraph()

		err := dg.ExecuteTargetKAndItsDeps("run")

		assert.ErrorIs(t, err, ErrTargetDoesnotExist, "fail to detect that target doesn't exist")
	})

	t.Run("Target has no commands", func(t *testing.T) {
		dg := NewDependencyGraph()
		target := "run"

		adjList := make(Graph)

		adjList[target] = []string{}
		dg.SetAdjacencyList(adjList)

		targetCommands := make(CommandMap)

		targetCommands[target] = []string{}
		dg.SetTargetToCommands(targetCommands)

		err := dg.ExecuteTargetKAndItsDeps(target)

		assert.ErrorIs(t, err, ErrTargetHasNoCommands)

	})

	// t.Run("valid target and command with @")
}
