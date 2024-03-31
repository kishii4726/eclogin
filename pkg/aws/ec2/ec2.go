package ec2

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func GetInstancesMap(c *ec2.Client) map[string]string {
	resp, err := c.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Fatalf("DescribeInstances failed %v\n", err)
	}

	ec2_reservations := resp.Reservations
	if len(ec2_reservations) == 0 {
		log.Fatalf("EC2 Instance does not exist")
	}

	instanceid_and_name_map := map[string]string{}
	for _, r := range ec2_reservations {
		var nameTag string
		for _, ins := range r.Instances {
			for _, tag := range ins.Tags {
				if aws.ToString(tag.Key) == "Name" {
					nameTag = aws.ToString(tag.Value)
					break
				}
			}
			if nameTag == "" {
				nameTag = "No Name Tag"
			}
			instanceid_and_name_map[fmt.Sprintf("%s(%s)", nameTag, *ins.InstanceId)] = *ins.InstanceId
		}
	}

	return instanceid_and_name_map
}

func GetInstances(m map[string]string) []string {
	var instances []string
	for k, _ := range m {
		instances = append(instances, k)
	}

	return instances
}
