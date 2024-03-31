package cmd

import (
	"encoding/json"
	"fmt"
	"smsh/pkg/aws/config"
	"smsh/pkg/aws/ecs"
	"smsh/pkg/aws/session"
	"smsh/pkg/prompt"

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

		cluster := prompt.GetUserSelectionFromList("Select ECS Cluster", ecs.GetClusters(client))
		service := prompt.GetUserSelectionFromList("Select ECS Service", ecs.GetServices(client, cluster))
		task_id := prompt.GetUserSelectionFromList("Select ECS Task Id", ecs.GetTaskIds(client, cluster, service))

		container_and_runtimeids := ecs.GetContainerAndRuntimeIDs(client, cluster, task_id)

		containers := ecs.GetContainers(container_and_runtimeids)
		container := prompt.GetUserSelectionFromList("Select ECS Container", containers)
		runtime_id := container_and_runtimeids[container]

		shell := prompt.GetUserSelectionFromList("Select Shell", []string{"sh", "bash"})

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
