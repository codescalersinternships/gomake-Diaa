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
		return "", err
	}

	cmd := exec.Command(path, strings.Join(cmdWords[1:], " "))

	if errors.Is(cmd.Err, exec.ErrDot) {
		return "", cmd.Err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err = cmd.Start(); err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		cmdOutput += fmt.Sprint(scanner.Text(), "\n")
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return cmdOutput, nil
}
