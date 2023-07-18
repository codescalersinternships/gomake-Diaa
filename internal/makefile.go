package makefile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var errInvalidFormat = errors.New("invalid format")

// ReadMakefile reads a Makefile from the specified file path and returns its dependency graph and command map.
func ReadMakefile(filePath string) (graph, commandMap, error) {

	file, err := os.Open(filePath)

	if err != nil {
		return nil, nil, err
	}

	defer file.Close()

	return parseMakefile(file)

}

func parseMakefile(r io.Reader) (graph, commandMap, error) {

	adjacencyList := make(graph)
	targetsCommands := make(commandMap)
	scanner := bufio.NewScanner(r)

	lineNum := 1
	currentTarget := ""
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		line = strings.Trim(line, " ")

		if isEmpty(line) || isComment(line) {
			continue
		}

		if isTarget(line) {

			target, deps := extractTargetAndDeps(line)

			if _, ok := adjacencyList[target]; ok {
				fmt.Printf("Warning: overriding recipe for target '%s'\n", target)
			}
			// make the graph
			currentTarget = target
			adjacencyList[currentTarget] = make([]string, 0)
			targetsCommands[currentTarget] = make([]string, 0)

			adjacencyList[currentTarget] = deps

			continue
		}

		isCommand := strings.HasPrefix(line, "\t")
		//command belongs to target
		if isCommand && currentTarget != "" {
			command := strings.TrimSpace(line)
			targetsCommands[currentTarget] = append(targetsCommands[currentTarget], command)
			continue
		}

		isGlobalCommand := isCommand && currentTarget == ""

		if isGlobalCommand {
			return nil, nil, fmt.Errorf("%w: global command at line %d", errInvalidFormat, lineNum)
		}
		// not comment. not command. not target. then invalid format
		return nil, nil, fmt.Errorf("%w at line %d", errInvalidFormat, lineNum)

	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return adjacencyList, targetsCommands, nil

}

func isTarget(line string) bool {
	targetAndDeps := strings.Split(line, ":")
	target, deps := targetAndDeps[0], targetAndDeps[1:]

	hasNoLeadingTab := !strings.HasPrefix(line, "\t")
	notEmptyTarget := len(target) > 0
	isOneTarget := !strings.Contains(target, " ")

	// invalid like target: dep1 :dep2
	isValidFormat := len(deps) == 1

	return hasNoLeadingTab && notEmptyTarget && isOneTarget && isValidFormat

}

func extractTargetAndDeps(line string) (string, []string) {

	targetAndDeps := strings.Split(line, ":")
	target, deps := targetAndDeps[0], targetAndDeps[1:]

	depsString := strings.TrimSpace(deps[0])
	deps = deps[:0]
	if len(depsString) > 0 {
		deps = strings.Split(depsString, " ")
	}
	return target, deps

}

func isEmpty(line string) bool {
	return len(strings.TrimSpace(line)) == 0
}

func isComment(line string) bool {
	line = strings.TrimSpace(line)

	return line[0] == '#'
}
