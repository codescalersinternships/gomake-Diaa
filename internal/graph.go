package makefile

import (
	"errors"
	"fmt"
	"strings"
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
	return &(DependencyGraph{adjacencyList: make(Graph),
		targetToCommands: make(CommandMap)})
}

func (d *DependencyGraph) SetAdjacencyList(adjList Graph) {
	d.adjacencyList = adjList
}

func (d *DependencyGraph) SetTargetToCommands(targCommands CommandMap) {
	d.targetToCommands = targCommands
}

func (d *DependencyGraph) checkCircularDependency() error {

	visited := make(map[string]bool)
	pathNodes := make(map[string]bool)

	for node := range d.adjacencyList {
		if !visited[node] {
			if err := d.checkCyclicPath(node, visited, pathNodes); err != nil {
				return err
			}
		}

	}
	return nil
}
func (d *DependencyGraph) checkCyclicPath(node string, visited, pathNodes map[string]bool) error {

	visited[node] = true
	pathNodes[node] = true

	for _, child := range d.adjacencyList[node] {
		if !visited[child] {
			return d.checkCyclicPath(child, visited, pathNodes)
		} else if pathNodes[child] {
			// cycle detected
			return fmt.Errorf("%w between '%s' -> '%s'", ErrCycleDetected, node, child)
		}
	}

	pathNodes[node] = false

	return nil
}

func (d *DependencyGraph) checkMissingDependencies() []string {
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

	if err := d.checkCircularDependency(); err != nil {
		return fmt.Errorf("cycle error: %w", err)
	}

	if misDep := d.checkMissingDependencies(); len(misDep) != 0 {
		return fmt.Errorf("missing dependencies: '%s'", strings.Join(misDep, ", "))
	}

	if _,ok := d.targetToCommands[target]; !ok {
		return fmt.Errorf("%w: 'target '%s' does not exist'", ErrTargetDoesnotExist, target)
	}

	visited := make(map[string]bool)

	_, err := d.executeTasksInDependencyOrder(target, visited)
	return err

}

func (d *DependencyGraph) executeTasksInDependencyOrder(target string, visited map[string]bool) (string, error) {

	visited[target] = true
	finalOutput := ""

	for _, child := range d.adjacencyList[target] {
		if !visited[child] {
			cmdOutput, err := d.executeTasksInDependencyOrder(child, visited)
			if err != nil {
				return "", err
			}
			finalOutput += cmdOutput
		}
	}

	// Exec commands of the leaf target
	cmdOutput, err := d.executeCommandsForTargetK(target)
	finalOutput += cmdOutput

	return finalOutput, err

}

func (d *DependencyGraph) executeCommandsForTargetK(target string) (string, error) {
	commands := d.targetToCommands[target]

	if len(commands) == 0 {
		return "", fmt.Errorf("%w for %s", ErrTargetHasNoCommands, target)
	}

	finalOutput := ""
	for _, command := range commands {

		commandOutput, err := execCommand(command)
		if err != nil {
			return "", err
		}
		finalOutput += commandOutput

	}
	return finalOutput, nil
}
