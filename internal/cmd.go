package makefile

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func execCommand(command string)  error{


	if command[0] == '@' {
		command = command[1:]
	}else{
		fmt.Println(command)
	}
	cmdWords := strings.Split(command, " ")

	cmd := exec.Command(cmdWords[0], strings.Join(cmdWords[1:], " "))

	if errors.Is(cmd.Err, exec.ErrDot) {
		return  cmd.Err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return  fmt.Errorf("pipe error: %w", err)
	}

	if err = cmd.Start(); err != nil {
		return  fmt.Errorf("error starting command: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("scanner err: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return  fmt.Errorf("error waiting the cmd execution: %w", err)
	}
	return nil
}
