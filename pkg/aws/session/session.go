package session

import (
	"os"
	"os/exec"
	"os/signal"
)

func ExecCommand(sessJson []byte, inputJson []byte, region string) {
	exec_cmd := exec.Command(
		"session-manager-plugin",
		string(sessJson),
		region,
		"StartSession",
		"",
		string(inputJson),
		"https://ssm."+region+".amazonaws.com",
	)
	signal.Ignore(os.Interrupt)
	defer signal.Reset(os.Interrupt)
	exec_cmd.Stdout = os.Stdout
	exec_cmd.Stdin = os.Stdin
	exec_cmd.Stderr = os.Stderr
	exec_cmd.Run()
}
