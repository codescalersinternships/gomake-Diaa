package makefile

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func execCommand(command string) (string, error) {

	cmdOutput := ""

	execQuietly := false
	if command[0] == '@' {
		execQuietly = true
		command = command[1:]
	}
	if !execQuietly {
		fmt.Println(command)
		cmdOutput += fmt.Sprint(command, "\n")
	}
	cmdWords := strings.Split(command, " ")
	binary := cmdWords[0]

	path, err := exec.LookPath(binary)

	if err != nil {
		return "", fmt.Errorf("binary does not exist:%w", err)
	}

	cmd := exec.Command(path, strings.Join(cmdWords[1:], " "))

	if errors.Is(cmd.Err, exec.ErrDot) {
		return "", cmd.Err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("pipe error: %w", err)
	}

	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting command: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		cmdOutput += fmt.Sprint(scanner.Text(), "\n")
	}
	if err = scanner.Err(); err != nil {
		return "", fmt.Errorf("scanner err: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("error waiting the cmd execution: %w", err)
	}
	return cmdOutput, nil
}
