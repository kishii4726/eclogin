package ecs

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func GetClusters(c *ecs.Client) []string {
	resp, err := c.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		log.Fatalf("ListClusters failed %v\n", err)
	}
	ecs_cluster_arns := resp.ClusterArns
	if len(ecs_cluster_arns) == 0 {
		log.Fatalf("Cluster does not exist")
	}
	ecs_clusters := []string{}
	for _, v := range ecs_cluster_arns {
		ecs_clusters = append(ecs_clusters, strings.Split(v, "/")[1])
	}

	return ecs_clusters
}

func GetServices(client *ecs.Client, ecs_cluster string) []string {
	resp, err := client.ListServices(context.TODO(), &ecs.ListServicesInput{
		Cluster: aws.String(ecs_cluster),
	})
	if err != nil {
		log.Fatalf("ListServices failed %v\n", err)
	}
	ecs_service_arns := resp.ServiceArns
	if len(ecs_service_arns) == 0 {
		log.Fatalf("Service does not exist")
	}
	ecs_services := []string{}
	for _, v := range ecs_service_arns {
		ecs_services = append(ecs_services, strings.Split(v, "/")[2])
	}

	return ecs_services
}

func GetTaskIds(client *ecs.Client, ecs_cluster string, ecs_service string) []string {
	resp, err := client.ListTasks(context.TODO(), &ecs.ListTasksInput{
		Cluster:     aws.String(ecs_cluster),
		ServiceName: aws.String(ecs_service),
	})
	if err != nil {
		log.Fatalf("ListTasks failed %v\n", err)
	}
	ecs_task_arns := resp.TaskArns
	if len(ecs_task_arns) == 0 {
		log.Fatalf("Task does not exist")
	}
	ecs_task_ids := []string{}
	for _, v := range ecs_task_arns {
		ecs_task_ids = append(ecs_task_ids, strings.Split(v, "/")[2])
	}

	return ecs_task_ids
}

func GetContainerAndRuntimeIDs(client *ecs.Client, ecs_cluster string, ecs_task_id string) map[string]string {
	describe_tasks, _ := client.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
		Tasks:   []string{ecs_task_id},
		Cluster: aws.String(ecs_cluster),
	})
	var container_and_runtimeids map[string]string
	for _, v := range describe_tasks.Tasks[0].Containers {
		container_and_runtimeids[*v.Name] = strings.Split(*v.RuntimeId, "-")[0]
	}

	return container_and_runtimeids
}

func GetContainers(container_and_runtimeids map[string]string) []string {
	var containers []string
	for c := range container_and_runtimeids {
		containers = append(containers, c)
	}

	return containers
}

func GetExecuteCommandOutput(client *ecs.Client, shell string, task_id string, cluster string, container string) *ecs.ExecuteCommandOutput {
	out, _ := client.ExecuteCommand(context.TODO(), &ecs.ExecuteCommandInput{
		// out, err := client.ExecuteCommand(context.TODO(), &ecs.ExecuteCommandInput{
		Command:     aws.String(shell),
		Interactive: true,
		Task:        aws.String(task_id),
		Cluster:     aws.String(cluster),
		Container:   aws.String(container),
	})

	return out
}
