package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	makefile "github.com/codescalersinternships/gomake-Diaa/internal"
)


func main() {
	filePath, target, err := parseInputCommand()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing command: %v\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	adjacencyList, targetToCommands, err := makefile.ReadMakefile(filePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	depGraph := makefile.NewDependencyGraph(adjacencyList,targetToCommands)

	err = depGraph.ExecuteTargetAndItsDeps(target)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// ParseInputCommand reads command line arguments and returns the file path and target to be executed.
// If the file path is not specified. its default value is 'Makefile'
func parseInputCommand() (string, string, error) {
	filePath := flag.String("f", "Makefile", "name of the file to be explored")
	target := flag.String("t", "", "target you want to execute")
	flag.Parse()

	if *target == "" {
		return "", "", errors.New("please specify a single target")
	}

	return *filePath, *target, nil
}
