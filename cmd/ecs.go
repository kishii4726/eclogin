package cmd

import (
	"context"
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

type ECSClientInterface interface {
	ListClusters(ctx context.Context, params *aws_ecs.ListClustersInput, optFns ...func(*aws_ecs.Options)) (*aws_ecs.ListClustersOutput, error)
	ListServices(ctx context.Context, params *aws_ecs.ListServicesInput, optFns ...func(*aws_ecs.Options)) (*aws_ecs.ListServicesOutput, error)
	ListTasks(ctx context.Context, params *aws_ecs.ListTasksInput, optFns ...func(*aws_ecs.Options)) (*aws_ecs.ListTasksOutput, error)
	DescribeTasks(ctx context.Context, params *aws_ecs.DescribeTasksInput, optFns ...func(*aws_ecs.Options)) (*aws_ecs.DescribeTasksOutput, error)
	ExecuteCommand(ctx context.Context, params *aws_ecs.ExecuteCommandInput, optFns ...func(*aws_ecs.Options)) (*aws_ecs.ExecuteCommandOutput, error)
}

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
	prompter := prompt.NewUIPrompter()
	region := prompt.GetFlagOrInput(cmd, "region", "Please enter AWS region", defaultRegion, prompter)
	profile := prompt.GetFlagOrInput(cmd, "profile", "Please enter AWS profile (optional)", "", prompter)

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

	container, runtimeID := selectContainer(cmd, containerInfo)
	shell := prompt.GetFlagOrSelect(cmd, "shell", "Select Shell", availableShells, prompter)

	printEcloginEcsWithOptionCommand(cluster, taskID, container, shell, region, profile)
	printAwsCliEcsCommand(cluster, taskID, container, shell, region, profile)

	if err := executeContainerSession(ecsClient, shell, taskID, cluster, container, runtimeID, region); err != nil {
		log.Fatalf("Failed to execute container session: %v", err)
	}
}

func getECSCluster(cmd *cobra.Command, client ECSClientInterface) (string, error) {
	clusters, err := ecs.ListClusters(client)
	if err != nil {
		return "", err
	}
	return prompt.GetFlagOrSelect(cmd, "cluster", "Select ECS Cluster", clusters, prompt.NewUIPrompter()), nil
}

func getECSService(cmd *cobra.Command, client ECSClientInterface, cluster string) (string, error) {
	services, err := ecs.ListServices(client, cluster)
	if err != nil {
		return "", err
	}
	return prompt.GetFlagOrSelect(cmd, "service", "Select ECS Service", services, prompt.NewUIPrompter()), nil
}

func getECSTaskID(cmd *cobra.Command, client ECSClientInterface, cluster, service string) (string, error) {
	taskIDs, err := ecs.ListTaskIDs(client, cluster, service)
	if err != nil {
		return "", err
	}
	return prompt.GetFlagOrSelect(cmd, "task-id", "Select ECS Task ID", taskIDs, prompt.NewUIPrompter()), nil
}

func selectContainer(cmd *cobra.Command, containerInfo map[string]string) (string, string) {
	containers := ecs.ListContainerNames(containerInfo)
	container := prompt.GetFlagOrSelect(cmd, "container", "Select ECS Container", containers, prompt.NewUIPrompter())
	return container, containerInfo[container]
}

func printEcloginEcsWithOptionCommand(cluster, taskID, container, shell, region, profile string) {
	if profile != "" {
		fmt.Printf(`eclogin equivalent command:
eclogin ecs --cluster %s --task-id %s --container %s --shell %s --region %s
`,
			cluster, taskID, container, shell, region)
	} else {
		fmt.Printf(`eclogin equivalent command:
eclogin ecs --cluster %s --task-id %s --container %s --shell %s --region %s --profile %s
`,
			cluster, taskID, container, shell, region, profile)
	}
}

func printAwsCliEcsCommand(cluster, taskID, container, shell, region, profile string) {
	if profile != "" {
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
	} else {
		fmt.Printf(`If you are using awscli, please copy the following:
aws ecs execute-command \
	--cluster %s \
	--task %s \
	--container %s \
	--interactive \
	--command %s \
	--region %s \
	--profile %s
`,
			cluster, taskID, container, shell, region, profile)
	}
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
