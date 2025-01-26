package cmd

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func chooseValueFromPromptItems(s string, l []string) string {
	prompt := promptui.Select{
		Label: s,
		Items: l,
	}
	_, v, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return v
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Execute commands inside local Docker containers.",
	Long: `The "local" command allows you to interact with and execute commands inside local Docker containers.
This command provides an interactive prompt to select a running or stopped container from your local Docker environment and choose a shell (such as /bin/sh or /bin/bash) to execute within the container.
The tool then establishes an interactive terminal session, allowing you to run commands directly inside the selected container.
`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Create a new Docker client with environment variables
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}

		// Update the client with the server's API version
		cli.NegotiateAPIVersion(ctx)

		containers, err := cli.ContainerList(ctx, container.ListOptions{Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "status",
			Value: "running",
		})})
		if err != nil {
			panic(err)
		}
		containers_map := map[string]string{}
		var containers_name_slice []string
		for _, container := range containers {
			for _, a := range container.Names {
				containers_map[strings.TrimLeft(a, "/")+"("+container.Image+")"] = container.ID
				containers_name_slice = append(containers_name_slice, strings.TrimLeft(a, "/")+"("+container.Image+")")
			}
		}

		selected_container := chooseValueFromPromptItems("Select Container", containers_name_slice)

		selected_shell := chooseValueFromPromptItems("Select Shell", []string{"/bin/sh", "/bin/bash"})

		execOpts := container.ExecOptions{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
			Cmd:          []string{selected_shell},
		}

		resp, err := cli.ContainerExecCreate(context.Background(), containers_map[selected_container], execOpts)
		if err != nil {
			panic(err)
		}

		execAttachResp, err := cli.ContainerExecAttach(context.Background(), resp.ID, container.ExecStartOptions{})
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := execAttachResp.Conn.Close(); err != nil {
				log.Panic(err)
			}
			log.Println("connection closed")
		}()

		fd := int(os.Stdin.Fd())
		if term.IsTerminal(fd) {
			state, err := term.MakeRaw(fd)
			if err != nil {
				log.Panic(err)
			}
			defer term.Restore(fd, state)
		}

		go io.Copy(execAttachResp.Conn, os.Stdin)
		stdcopy.StdCopy(os.Stdout, os.Stderr, execAttachResp.Reader)
	},
}

func init() {
	rootCmd.AddCommand(localCmd)
}
