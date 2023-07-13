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
$ ./gomake [options] [target]
```

### Options

- -f FILE: Specify the path to the Makefile. If not specified, the default path Makefile will be used.

## Features
- **Target ordering:** GoMake orders the targets based on their dependencies, ensuring that each target is executed only after its dependencies are completed.
- **Circular dependency detection:** GoMake checks for circular dependencies between targets and returns an error if a circular dependency is detected.

- **Missing dependency detection:** GoMake checks for missing dependencies and returns an error if a target has a dependency that is not defined in the Makefile.

- **Command execution:** GoMake executes the commands associated with each target in order, ensuring that the commands are executed only after their dependencies are completed.

- **Error handling** GoMake returns detailed error messages for circular dependencies, missing dependencies, and targets with no commands.
