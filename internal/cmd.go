package makefile

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func execCommand(command string) error {

	if command[0] == '@' {
		command = command[1:]
	} else {
		fmt.Println(command)
	}
	cmdWords := strings.Split(command, " ")


	cmd := exec.Command(cmdWords[0], cmdWords[1:]...)

	if errors.Is(cmd.Err, exec.ErrDot) {
		return cmd.Err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}
