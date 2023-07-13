# GoMake

GoMake is a simple implementation of the make command in Go. It reads a Makefile and executes the specified target and its dependencies.

## Installation

1. Clone the repo:

```git
$ git clone https://github.com/codescalersinternships/gomake-Diaa.git
```

2. Navigate to the repo directory.

```go
$ cd gomake-Diaa
```

3. Install dependencies

```go
$ go get -d ./...
```

4. Build the package:

```go
$ go build -o "bin/gomake" main.go
```

5. Navigate to bin

```git
$ cd bin
```

## Usage

```go
$ ./gomake [options] -t [target]
```

### Options

- -f FILE: Specify the path to the Makefile. If not specified, the default path Makefile will be used.
- target: target you want to execute

## Features
- **Target ordering:** GoMake orders the targets based on their dependencies, ensuring that each target is executed only after its dependencies are completed.
- **Circular dependency detection:** GoMake checks for circular dependencies between targets and returns an error if a circular dependency is detected.

- **Missing dependency detection:** GoMake checks for missing dependencies and returns an error if a target has a dependency that is not defined in the Makefile.

- **Command execution:** GoMake executes the commands associated with each target in order, ensuring that the commands are executed only after their dependencies are completed.

- **Error handling** GoMake returns detailed error messages for circular dependencies, missing dependencies, and targets with no commands.

## Format
GoMake checks a set of rules, each of which specifies a target and its dependencies, along with the commands needed to build the target. Rules have the following format:
```make
target: [dependency ...]
[tab] command
```

- target: the name of the target to build.
- dependency: the name of a target that must be built before the current target can be built.
- command: the command to execute to build the target.
- The commands to build the target are indented with a single tab character
- global commands not allowed
- comments starts with '#'
- comments and empty lines are ignored

Here is an example of a simple Makefile:

```make
build: main.o lib.o
	go build -o myprogram main.o lib.o

main.o: main.go
	go build -o main.o -c main.go

lib.o: lib.go
	go build -o lib.o -c lib.go
```