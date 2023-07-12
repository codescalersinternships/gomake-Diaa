package main

import (
	"errors"
	"flag"
	"fmt"
	makefile "gomake/internal"
	"os"
	"strings"
)

var ErrMissingMakefileArg = fmt.Errorf("make: option requires an argument -- 'f'")

const HelpMessage = `Usage: make [options] [target] ...
Options:
  -f FILE`

func main() {
	filePath, target, err := ParseCommand()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing command: %v\n", err)
		fmt.Println(HelpMessage)
		os.Exit(1)
	}

	adjList, targToCmds, err := makefile.ReadMakefile(filePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	depGraph := makefile.NewDependencyGraph()

	depGraph.SetAdjacencyList(adjList)
	depGraph.SetTargetToCommands(targToCmds)

	// check cycles
	if err = depGraph.HasCircularDependency(); err != nil {

		fmt.Fprintf(os.Stderr, "cycle error: %v\n", err)
		os.Exit(1)
	}

	// check missing deps
	if misDep := depGraph.CheckMissingDependencies(); len(misDep) != 0 {
		fmt.Fprintf(os.Stderr, "error missing dependencies: %s\n", strings.Join(misDep, ", "))
		os.Exit(1)
	}

	err = depGraph.ExecuteTargetKAndItsDeps(target)

	if err != nil {
		fmt.Fprintf(os.Stderr, "execution error: %v\n", err)
		os.Exit(1)
	}

}

func ParseCommand() (string, string, error) {
	filePath := flag.String("f", "", "name of the file to be explored")

	flag.Parse()

	if len(flag.Args()) != 1 {
		return "", "", errors.New("Please specify a single target.")
	}
	target := flag.Args()[0]

	if *filePath == "" {
		*filePath = "Makefile"
	}

	return *filePath, target, nil
}

func SearchArray[T comparable](arr []T, target T) (bool, int) {
	for idx, value := range arr {
		if value == target {
			return true, idx
		}
	}
	return false, -1
}

// always has the filename
// detect cycles always , crash if there is any cycle
// one target
