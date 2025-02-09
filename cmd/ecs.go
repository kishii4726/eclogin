package cmd

import (
	"eclogin/pkg/aws/config"
	"eclogin/pkg/aws/ecs"
	"eclogin/pkg/aws/session"
	"eclogin/pkg/prompt"
	"encoding/json"
	"fmt"
	"log"

	aws_ecs "github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
)

const (
	defaultRegion = "ap-northeast-1"
	targetFormat  = "ecs:%s_%s_%s"
)

var (
	availableShells = []string{"sh", "bash"}
)

var ecsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "Start an interactive session with an ECS container using ECS Exec",
	Long: `The ecs command allows you to start an interactive session with an ECS container
using ECS Exec. You can select a cluster, service, task, and container,
and establish a session to manage it remotely.`,
	Run: runECSCommand,
}

func runECSCommand(cmd *cobra.Command, _ []string) {
	region := prompt.GetFlagOrInput(cmd, "region", "Please enter AWS region", defaultRegion)
	profile := prompt.GetFlagOrInput(cmd, "profile", "Please enter AWS profile (optional)", "")

	cfg, err := config.LoadConfig(region, profile)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	ecsClient := aws_ecs.NewFromConfig(cfg)

	cluster, err := getECSCluster(cmd, ecsClient)
	if err != nil {
		log.Fatalf("Failed to get ECS cluster: %v", err)
	}

	service, err := getECSService(cmd, ecsClient, cluster)
	if err != nil {
		log.Fatalf("Failed to get ECS service: %v", err)
	}

	taskID, err := getECSTaskID(cmd, ecsClient, cluster, service)
	if err != nil {
		log.Fatalf("Failed to get ECS task ID: %v", err)
	}

	containerInfo, err := ecs.GetContainerInfo(ecsClient, cluster, taskID)
	if err != nil {
		log.Fatalf("Failed to get container information: %v", err)
	}

	container, runtimeID := selectContainer(containerInfo)
	shell := prompt.PromptSelect("Select Shell", availableShells)

	printAwsCliEcsCommand(cluster, taskID, container, shell, region)

	if err := executeContainerSession(ecsClient, shell, taskID, cluster, container, runtimeID, region); err != nil {
		log.Fatalf("Failed to execute container session: %v", err)
	}
}

func getECSCluster(cmd *cobra.Command, client *aws_ecs.Client) (string, error) {
	clusters, err := ecs.ListClusters(client)
	if err != nil {
		return "", err
	}
	return prompt.GetFlagOrSelect(cmd, "cluster", "Select ECS Cluster", clusters), nil
}

func getECSService(cmd *cobra.Command, client *aws_ecs.Client, cluster string) (string, error) {
	services, err := ecs.ListServices(client, cluster)
	if err != nil {
		return "", err
	}
	return prompt.GetFlagOrSelect(cmd, "service", "Select ECS Service", services), nil
}

func getECSTaskID(cmd *cobra.Command, client *aws_ecs.Client, cluster, service string) (string, error) {
	taskIDs, err := ecs.ListTaskIDs(client, cluster, service)
	if err != nil {
		return "", err
	}
	return prompt.GetFlagOrSelect(cmd, "task-id", "Select ECS Task ID", taskIDs), nil
}

func selectContainer(containerInfo map[string]string) (string, string) {
	containers := ecs.ListContainerNames(containerInfo)
	container := prompt.PromptSelect("Select ECS Container", containers)
	return container, containerInfo[container]
}

func printAwsCliEcsCommand(cluster, taskID, container, shell, region string) {
	fmt.Printf(`If you are using awscli, please copy the following:
aws ecs execute-command \
	--cluster %s \
	--task %s \
	--container %s \
	--interactive \
	--command %s \
	--region %s
`,
		cluster, taskID, container, shell, region)
}

func executeContainerSession(client *aws_ecs.Client, shell, taskID, cluster, container, runtimeID, region string) error {
	out, err := ecs.ExecuteContainerCommand(client, shell, taskID, cluster, container)
	if err != nil {
		return fmt.Errorf("execute command failed: %w", err)
	}

	sessionJSON, err := json.Marshal(out.Session)
	if err != nil {
		return fmt.Errorf("marshal session failed: %w", err)
	}

	target := fmt.Sprintf(targetFormat, cluster, taskID, runtimeID)
	input := ssm.StartSessionInput{
		Target: &target,
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshal input failed: %w", err)
	}

	return session.StartSession(sessionJSON, inputJSON, region)
}

func init() {
	rootCmd.AddCommand(ecsCmd)
}
