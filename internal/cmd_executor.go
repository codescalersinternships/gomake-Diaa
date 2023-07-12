package makefile

import "fmt"

func CMD_Exec(command string) error {
	// gomakeExec, err := exec.LookPath("gomake")

	// if err!=nil{
	// 	return err
	// }

	// cmd:=exec.Command(gomakeExec, command)

	fmt.Println(command)
	return nil
}
