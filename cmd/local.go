package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	defaultShellBash = "/bin/bash"
	defaultShellSh   = "/bin/sh"
)

type dockerExecutor struct {
	client *client.Client
	ctx    context.Context
}

func newDockerExecutor(ctx context.Context) (*dockerExecutor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)
	return &dockerExecutor{client: cli, ctx: ctx}, nil
}

func (d *dockerExecutor) getRunningContainers() (map[string]string, []string, error) {
	containerFilter := filters.NewArgs(filters.KeyValuePair{
		Key:   "status",
		Value: "running",
	})

	containers, err := d.client.ContainerList(d.ctx, container.ListOptions{Filters: containerFilter})
	if err != nil {
		return nil, nil, err
	}

	containerMap := make(map[string]string)
	var containerNames []string

	for _, container := range containers {
		for _, name := range container.Names {
			displayName := formatContainerName(name, container.Image)
			containerMap[displayName] = container.ID
			containerNames = append(containerNames, displayName)
		}
	}

	return containerMap, containerNames, nil
}

func formatContainerName(name, image string) string {
	return strings.TrimLeft(name, "/") + "(" + image + ")"
}

func selectOption(prompt string, options []string) (string, error) {
	selector := promptui.Select{
		Label: prompt,
		Items: options,
	}
	_, result, err := selector.Run()
	return result, err
}

func (d *dockerExecutor) executeInContainer(containerID, shell string) error {
	execConfig := container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{shell},
	}

	execResp, err := d.client.ContainerExecCreate(d.ctx, containerID, execConfig)
	if err != nil {
		return err
	}

	execConn, err := d.client.ContainerExecAttach(d.ctx, execResp.ID, container.ExecStartOptions{})
	if err != nil {
		return err
	}
	defer execConn.Close()

	return setupTerminal(execConn)
}
func setupTerminal(execConn types.HijackedResponse) error {
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		state, err := term.MakeRaw(fd)
		if err != nil {
			return err
		}
		defer term.Restore(fd, state)
	}

	go io.Copy(execConn.Conn, os.Stdin)
	stdcopy.StdCopy(os.Stdout, os.Stderr, execConn.Reader)
	return nil
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Execute commands inside local Docker containers.",
	Long: `The "local" command allows you to interact with and execute commands inside local Docker containers.
This command provides an interactive prompt to select a running or stopped container from your local Docker environment and choose a shell (such as /bin/sh or /bin/bash) to execute within the container.
The tool then establishes an interactive terminal session, allowing you to run commands directly inside the selected container.
`,
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx := context.Background()

		executor, err := newDockerExecutor(ctx)
		if err != nil {
			return err
		}

		containerMap, containerNames, err := executor.getRunningContainers()
		if err != nil {
			return err
		}

		if len(containerNames) == 0 {
			return fmt.Errorf("no running containers found")
		}

		selectedContainer, err := selectOption("Select Container", containerNames)
		if err != nil {
			return err
		}

		selectedShell, err := selectOption("Select Shell", []string{defaultShellSh, defaultShellBash})
		if err != nil {
			return err
		}

		return executor.executeInContainer(containerMap[selectedContainer], selectedShell)
	},
}

func init() {
	rootCmd.AddCommand(localCmd)
}
