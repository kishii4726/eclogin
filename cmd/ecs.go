package cmd

import (
	"eclogin/pkg/aws/config"
	"eclogin/pkg/aws/ecs"
	"eclogin/pkg/aws/session"
	"eclogin/pkg/prompt"
	"encoding/json"
	"fmt"

	aws_ecs "github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
)

var ecsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "Start an interactive session with an ECS container using ECS Exec",
	Long: `The ecs command allows you to start an interactive session with an ECS container
using ECS Exec. You can select a cluster, service, task, and container,
and establish a session to manage it remotely.`,
	Run: func(cmd *cobra.Command, args []string) {
		region := prompt.GetFlag(cmd, "region", "Please enter aws region(Default: ap-northeast-1)", "ap-northeast-1")
		profile := prompt.GetFlag(cmd, "profile", "Please enter aws profile(If empty, default settings are loaded)", "")

		cfg := config.LoadConfig(region, profile)
		client := aws_ecs.NewFromConfig(cfg)

		cluster := prompt.GetFlagOrPrompt(cmd, "cluster", "Select ECS Cluster", func() []string { return ecs.GetCluster(client) })
		service := prompt.GetFlagOrPrompt(cmd, "service", "Select ECS Service", func() []string { return ecs.GetService(client, cluster) })
		task_id := prompt.GetFlagOrPrompt(cmd, "task-id", "Select ECS Task Id", func() []string { return ecs.GetTaskId(client, cluster, service) })

		container_and_runtimeids := ecs.GetContainerAndRuntimeIDs(client, cluster, task_id)

		containers := ecs.GetContainers(container_and_runtimeids)
		container := prompt.GetUserSelectionFromList("Select ECS Container", containers)
		runtime_id := container_and_runtimeids[container]

		shell := prompt.GetUserSelectionFromList("Select Shell", []string{"sh", "bash"})

		fmt.Printf(`If you are using awscli, please copy the following:
aws ecs execute-command \
	--cluster %s \
	--task %s \
	--container %s \
	--interactive \
	--command %s \
	--region %s
`,
			cluster, task_id, container, shell, region)

		out := ecs.GetExecuteCommandOutput(client, shell, task_id, cluster, container)
		sessJson, _ := json.Marshal(out.Session)
		target := fmt.Sprintf("ecs:%s_%s_%s", cluster, task_id, runtime_id)
		input := ssm.StartSessionInput{
			Target: &target,
		}
		inputJson, _ := json.Marshal(input)
		session.ExecCommand(sessJson, inputJson, region)
	},
}

func init() {
	rootCmd.AddCommand(ecsCmd)
}
