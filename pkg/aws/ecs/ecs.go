package ecs

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ECSClient interface {
	ListClusters(ctx context.Context, params *ecs.ListClustersInput, optFns ...func(*ecs.Options)) (*ecs.ListClustersOutput, error)
	ListServices(ctx context.Context, params *ecs.ListServicesInput, optFns ...func(*ecs.Options)) (*ecs.ListServicesOutput, error)
	ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error)
	DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)
	ExecuteCommand(ctx context.Context, params *ecs.ExecuteCommandInput, optFns ...func(*ecs.Options)) (*ecs.ExecuteCommandOutput, error)
}

func ListClusters(c ECSClient) ([]string, error) {
	resp, err := c.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	clusterARNs := resp.ClusterArns
	if len(clusterARNs) == 0 {
		return nil, fmt.Errorf("no clusters found")
	}

	clusters := make([]string, len(clusterARNs))
	for i, arn := range clusterARNs {
		clusters[i] = strings.Split(arn, "/")[1]
	}

	return clusters, nil
}

func ListServices(client ECSClient, clusterName string) ([]string, error) {
	resp, err := client.ListServices(context.TODO(), &ecs.ListServicesInput{
		Cluster:    aws.String(clusterName),
		MaxResults: aws.Int32(100),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	serviceARNs := resp.ServiceArns
	if len(serviceARNs) == 0 {
		return nil, fmt.Errorf("no services found in cluster %s", clusterName)
	}

	services := make([]string, len(serviceARNs))
	for i, arn := range serviceARNs {
		parts := strings.Split(arn, "/")
		services[i] = parts[len(parts)-1]
	}

	return services, nil
}

func ListTaskIDs(client ECSClient, clusterName, serviceName string) ([]string, error) {
	resp, err := client.ListTasks(context.TODO(), &ecs.ListTasksInput{
		Cluster:     aws.String(clusterName),
		ServiceName: aws.String(serviceName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	taskARNs := resp.TaskArns
	if len(taskARNs) == 0 {
		return nil, fmt.Errorf("no tasks found for service %s in cluster %s", serviceName, clusterName)
	}

	taskIDs := make([]string, len(taskARNs))
	for i, arn := range taskARNs {
		taskIDs[i] = strings.Split(arn, "/")[2]
	}

	return taskIDs, nil
}

func GetContainerInfo(client ECSClient, clusterName, taskID string) (map[string]string, error) {
	resp, err := client.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
		Tasks:   []string{taskID},
		Cluster: aws.String(clusterName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe tasks: %w", err)
	}

	containerInfo := make(map[string]string)
	for _, container := range resp.Tasks[0].Containers {
		containerInfo[*container.Name] = strings.Split(*container.RuntimeId, "-")[0]
	}

	return containerInfo, nil
}

func ListContainerNames(containerInfo map[string]string) []string {
	containers := make([]string, 0, len(containerInfo))
	for name := range containerInfo {
		containers = append(containers, name)
	}
	return containers
}

func ExecuteContainerCommand(client ECSClient, command, taskID, clusterName, containerName string) (*ecs.ExecuteCommandOutput, error) {
	output, err := client.ExecuteCommand(context.TODO(), &ecs.ExecuteCommandInput{
		Command:     aws.String(command),
		Interactive: true,
		Task:        aws.String(taskID),
		Cluster:     aws.String(clusterName),
		Container:   aws.String(containerName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	return output, nil
}
