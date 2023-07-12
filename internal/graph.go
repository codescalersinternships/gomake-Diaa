package makefile

import (
	"errors"
	"fmt"
)

var TargetDoesnotExistErr = errors.New("target error:")

type Graph = map[string][]string
type CommandMap = map[string][]string

type DependencyGraph struct {
	adjacencyList    map[string][]string
	targetToCommands map[string][]string
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


func (d *DependencyGraph) HasCircularDependency() error {

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
			return fmt.Errorf("circular dependency detected between '%s' -> '%s'", node, child)
		}
	}

	pathNodes[node] = false

	return nil
}

func (d *DependencyGraph) CheckMissingDependencies() []string {
	missingDeps := make([]string, 0)

	for target, deps := range d.adjacencyList {
		for _, dep := range deps {
			_, ok := d.adjacencyList[dep]
			if !ok {

				missingDeps = append(missingDeps, fmt.Sprintf("%s -> %s", target, dep))
			}
		}
	}
	return missingDeps
}

func (d *DependencyGraph) ExecuteTargetKAndItsDeps(target string) error {

	if _, ok := d.adjacencyList[target]; !ok {
		return fmt.Errorf("%v target %s does not exist", TargetDoesnotExistErr, target)
	}

	visited := make(map[string]bool)

	return d.ExecTasks(target, visited)

}

func (d *DependencyGraph) ExecTasks(target string, visited map[string]bool) error {

	visited[target] = true

	for _, child := range d.adjacencyList[target] {
		if !visited[child] {
			if err := d.ExecTasks(child, visited); err != nil {
				return err
			}
		}
	}

	// Exec commands of the leaf target
	return d.ExecTargetKCommands(target)

}

func (d *DependencyGraph) ExecTargetKCommands(target string) error {
	commands := d.targetToCommands[target]

	if len(commands) == 0 {
		return fmt.Errorf("gomake: Nothing to be done for %s", target)
	}

	for _, command := range commands {
		if command[0] == '@' {
			fmt.Println(command)
			command = command[1:]
		}

		err := CMD_Exec(command)
		if err != nil {
			return err
		}

	}
	return nil
}
