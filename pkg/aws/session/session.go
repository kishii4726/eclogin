package session

import (
	"os"
	"os/exec"
	"os/signal"
)

func StartSession(sessionData []byte, inputData []byte, region string) error {
	cmd := exec.Command(
		"session-manager-plugin",
		string(sessionData),
		region,
		"StartSession",
		"",
		string(inputData),
		"https://ssm."+region+".amazonaws.com",
	)
	signal.Ignore(os.Interrupt)
	defer signal.Reset(os.Interrupt)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
