package session

import (
	"os"
	"os/exec"
	"testing"
)

var execCommand func(command string, args ...string) *exec.Cmd

func TestExecCommand(t *testing.T) {
	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
		return cmd
	}

	sessJson := []byte(`{"sessionId":"test-session"}`)
	inputJson := []byte(`{"input":"test-input"}`)
	region := "ap-northeast-1"

	StartSession(sessJson, inputJson, region)
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}
