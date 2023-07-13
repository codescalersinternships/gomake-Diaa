package makefile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrInvalidFormat = errors.New("invalid format")

func ReadMakefile(filePath string) (Graph, CommandMap, error) {

	file, err := os.Open(filePath)

	if err != nil {
		return nil, nil, err
	}

	defer file.Close()

	return ParseMakefile(file)

}

func ParseMakefile(r io.Reader) (Graph, CommandMap, error) {

	adjList := make(Graph)
	targetsCommands := make(CommandMap)
	scanner := bufio.NewScanner(r)

	lineNum := 1
	currentTarget := ""
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if isEmpty(line) || isComment(line) {
			continue
		}

		target, deps := extractTargetAndDeps(line)
		isTarget := target != ""

		if isTarget {

			if _, ok := adjList[target]; ok {
				fmt.Printf("Warning: overriding recipe for target '%s'\n", target)
			}
			// make the graph
			currentTarget = target
			adjList[currentTarget] = make([]string, 0)
			targetsCommands[currentTarget] = make([]string, 0)

			adjList[currentTarget] = deps

			continue
		}

		command, isCommand := extractCommand(line)

		//command belongs to target
		if isCommand && currentTarget != "" {
			targetsCommands[currentTarget] = append(targetsCommands[currentTarget], command)
			continue
		}

		isGlobalCommand := isCommand && currentTarget == ""

		if isGlobalCommand {
			return nil, nil, fmt.Errorf("%w: global command at line %d", ErrInvalidFormat, lineNum)

		} else {
			// not comment. not command. not target. then invalid format
			return nil, nil, fmt.Errorf("%w at line %d", ErrInvalidFormat, lineNum)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return adjList, targetsCommands, nil

}

func extractTargetAndDeps(line string) (string, []string) {

	line = strings.Trim(line, " ")

	targetAndDeps := strings.Split(line, ":")
	target, deps := targetAndDeps[0], targetAndDeps[1:]

	hasNoLeadingTab := !strings.HasPrefix(line, "\t")
	notEmptyTarget := len(target) > 0
	isOneTarget := !strings.Contains(target, " ")

	// invalid like target: dep1 :dep2
	isValidFormat := len(deps) == 1

	if hasNoLeadingTab && notEmptyTarget && isOneTarget && isValidFormat {

		depsString := strings.TrimSpace(deps[0])
		deps = deps[:0]
		if len(depsString) > 0 {
			deps = strings.Split(depsString, " ")
		}
		return target, deps
	}
	return "", nil
}

func extractCommand(line string) (string, bool) {
	line = strings.Trim(line, " ")

	if strings.HasPrefix(line, "\t") {
		return strings.TrimSpace(line), true
	}
	return "", false
}

func isEmpty(line string) bool {
	return len(strings.TrimSpace(line)) == 0
}

func isComment(line string) bool {
	line = strings.TrimSpace(line)

	return line[0] == '#'
}
