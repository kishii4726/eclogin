package ec2

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func GetInstanceNameIDMap(client *ec2.Client) map[string]string {
	instances, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Fatalf("Failed to describe EC2 instances: %v", err)
	}

	reservations := instances.Reservations
	if len(reservations) == 0 {
		log.Fatalf("No EC2 instances found")
	}

	instanceMap := make(map[string]string)
	for _, reservation := range reservations {
		for _, instance := range reservation.Instances {
			name := getInstanceName(instance.Tags)
			displayName := fmt.Sprintf("%s(%s)", name, *instance.InstanceId)
			instanceMap[displayName] = *instance.InstanceId
		}
	}

	return instanceMap
}
func getInstanceName(tags []types.Tag) string {
	for _, tag := range tags {
		if aws.ToString(tag.Key) == "Name" {
			return aws.ToString(tag.Value)
		}
	}
	return "No Name Tag"
}

func GetInstanceDisplayNames(instanceMap map[string]string) []string {
	displayNames := make([]string, 0, len(instanceMap))
	for name := range instanceMap {
		displayNames = append(displayNames, name)
	}
	return displayNames
}
