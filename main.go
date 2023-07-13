package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	makefile "github.com/codescalersinternships/gomake-Diaa/internal"
)

var ErrMissingMakefileArg = fmt.Errorf("make: option requires an argument -- 'f'")

const HelpMessage = `Usage: make [options] [target] ...
Options:
  -f FILE`

func main() {
	filePath, target, err := ParseInputCommand()

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

	err = depGraph.ExecuteTargetKAndItsDeps(target)

	if err != nil {
		fmt.Fprintf(os.Stderr, "execution error: %v\n", err)
		os.Exit(1)
	}
}

func ParseInputCommand() (string, string, error) {
	filePath := flag.String("f", "", "name of the file to be explored")
	target := flag.String("t", "", "target you want to execute")
	flag.Parse()

	if *target == "" {
		return "", "", errors.New("please specify a single target")
	}

	if *filePath == "" {
		*filePath = "Makefile"
	}

	return *filePath, *target, nil
}
