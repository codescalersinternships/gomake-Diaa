package makefile

import (
	"errors"
	"fmt"
)

var (
	ErrTargetDoesnotExist  = errors.New("target error")
	ErrCycleDetected       = errors.New("circular dependency detected")
	ErrTargetHasNoCommands = errors.New("gomake: Nothing to be done")
)

type Graph = map[string][]string
type CommandMap = map[string][]string

type DependencyGraph struct {
	adjacencyList    Graph
	targetToCommands CommandMap
}

func NewDependencyGraph() *DependencyGraph {
	return &(DependencyGraph{adjacencyList: make(map[string][]string),
		targetToCommands: make(map[string][]string)})
}

func (d *DependencyGraph) SetAdjacencyList(adjList Graph) {
	d.adjacencyList = adjList
}

func (d *DependencyGraph) SetTargetToCommands(targCommands CommandMap) {
	d.targetToCommands = targCommands
}

func (d *DependencyGraph) CheckCircularDependency() error {

	visited := make(map[string]bool)
	pathNodes := make(map[string]bool)

	for node := range d.adjacencyList {
		if !visited[node] {
			if err := d.CheckCyclicPath(node, visited, pathNodes); err != nil {
				return err
			}
		}

	}
	return nil
}
func (d *DependencyGraph) CheckCyclicPath(node string, visited, pathNodes map[string]bool) error {

	visited[node] = true
	pathNodes[node] = true

	for _, child := range d.adjacencyList[node] {
		if !visited[child] {
			return d.CheckCyclicPath(child, visited, pathNodes)
		} else if pathNodes[child] {
			// cycle detected
			return fmt.Errorf("%w between '%s' -> '%s'", ErrCycleDetected, node, child)
		}
	}

	pathNodes[node] = false

	return nil
}

func (d *DependencyGraph) CheckMissingDependencies() []string {
	missingDeps := make([]string, 0)

	for _, deps := range d.adjacencyList {
		for _, dep := range deps {
			_, ok := d.adjacencyList[dep]
			if !ok {

				missingDeps = append(missingDeps, dep)
			}
		}
	}
	return missingDeps
}

func (d *DependencyGraph) ExecuteTargetKAndItsDeps(target string) error {

	if _, ok := d.adjacencyList[target]; !ok {
		return fmt.Errorf("%w: 'target '%s' does not exist'", ErrTargetDoesnotExist, target)
	}

	visited := make(map[string]bool)

	return d.executeTasksInDependencyOrder(target, visited)

}

func (d *DependencyGraph) executeTasksInDependencyOrder(target string, visited map[string]bool) error {

	visited[target] = true

	for _, child := range d.adjacencyList[target] {
		if !visited[child] {
			if err := d.executeTasksInDependencyOrder(child, visited); err != nil {
				return err
			}
		}
	}

	// Exec commands of the leaf target
	return d.executeCommandsForTargetK(target)

}

func (d *DependencyGraph) executeCommandsForTargetK(target string) error {
	commands := d.targetToCommands[target]

	if len(commands) == 0 {
		return fmt.Errorf("%w for %s", ErrTargetHasNoCommands, target)
	}

	for _, command := range commands {
		execQuietly := false
		if command[0] == '@' {
			execQuietly = true
			command = command[1:]
		}
		if !execQuietly {

			fmt.Println(command)
		}
		err := execCommand(command)
		if err != nil {
			return err
		}

	}
	return nil
}
