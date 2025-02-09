package ecs

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type mockECSClient struct{}

func (m *mockECSClient) ListClusters(ctx context.Context, params *ecs.ListClustersInput, optFns ...func(*ecs.Options)) (*ecs.ListClustersOutput, error) {
	return &ecs.ListClustersOutput{
		ClusterArns: []string{"arn:aws:ecs:region:account-id:cluster/test-cluster"},
	}, nil
}

func (m *mockECSClient) ListServices(ctx context.Context, params *ecs.ListServicesInput, optFns ...func(*ecs.Options)) (*ecs.ListServicesOutput, error) {
	return &ecs.ListServicesOutput{
		ServiceArns: []string{"arn:aws:ecs:region:account-id:service/cluster/test-service"},
	}, nil
}

func (m *mockECSClient) ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error) {
	return &ecs.ListTasksOutput{
		TaskArns: []string{"arn:aws:ecs:region:account-id:task/test-cluster/test-task"},
	}, nil
}

func (m *mockECSClient) DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error) {
	containerName := "test-container"
	runtimeId := "test-runtime-id-container"
	return &ecs.DescribeTasksOutput{
		Tasks: []types.Task{
			{
				Containers: []types.Container{
					{
						Name:      &containerName,
						RuntimeId: &runtimeId,
					},
				},
			},
		},
	}, nil
}

func (m *mockECSClient) ExecuteCommand(ctx context.Context, params *ecs.ExecuteCommandInput, optFns ...func(*ecs.Options)) (*ecs.ExecuteCommandOutput, error) {
	return &ecs.ExecuteCommandOutput{}, nil
}

func TestListClusters(t *testing.T) {
	client := &mockECSClient{}
	clusters, _ := ListClusters(client)
	if len(clusters) != 1 || clusters[0] != "test-cluster" {
		t.Errorf("expected test-cluster, got %v", clusters)
	}
}

func TestListServices(t *testing.T) {
	client := &mockECSClient{}
	services, _ := ListServices(client, "test-cluster")
	if len(services) != 1 || services[0] != "test-service" {
		t.Errorf("expected test-service, got %v", services)
	}
}

func TestListTaskIDs(t *testing.T) {
	client := &mockECSClient{}
	tasks, _ := ListTaskIDs(client, "test-cluster", "test-service")
	if len(tasks) != 1 || tasks[0] != "test-task" {
		t.Errorf("expected test-task, got %v", tasks)
	}
}

func TestGetContainerInfo(t *testing.T) {
	client := &mockECSClient{}
	containerAndRuntimeIDs, _ := GetContainerInfo(client, "test-cluster", "test-task")
	if len(containerAndRuntimeIDs) != 1 || containerAndRuntimeIDs["test-container"] != "test" {
		t.Errorf("expected test, got %v", containerAndRuntimeIDs["test-container"])
	}
}

func TestListContainerNames(t *testing.T) {
	containerAndRuntimeIDs := map[string]string{"test-container": "test-runtime"}
	containers := ListContainerNames(containerAndRuntimeIDs)
	if len(containers) != 1 || containers[0] != "test-container" {
		t.Errorf("expected test-container, got %v", containers)
	}
}

func TestExecuteContainerCommand(t *testing.T) {
	client := &mockECSClient{}
	output, _ := ExecuteContainerCommand(client, "/bin/bash", "test-task", "test-cluster", "test-container")
	if output == nil {
		t.Error("expected non-nil output")
	}
}
