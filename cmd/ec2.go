package cmd

import (
	"context"
	"encoding/json"
	"log"
	"eclogin/pkg/aws/config"
	"eclogin/pkg/aws/ec2"
	"eclogin/pkg/aws/session"
	"eclogin/pkg/prompt"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
)

var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Start an interactive session with an EC2 instance using AWS Systems Manager",
	Long: `The ec2 command allows you to start an interactive session with an EC2 instance
using AWS Systems Manager. You can select an instance from the list of running instances
and establish a session to manage it remotely.`,
	Run: func(cmd *cobra.Command, args []string) {
		region := prompt.GetFlag(cmd, "region", "Please enter aws region(Default: ap-northeast-1)", "ap-northeast-1")
		profile := prompt.GetFlag(cmd, "profile", "Please enter aws profile(If empty, default settings are loaded)", "")

		cfg := config.LoadConfig(region, profile)
		client := aws_ec2.NewFromConfig(cfg)
		instanceid_and_name_map := ec2.GetInstancesMap(client)

		instances := ec2.GetInstances(instanceid_and_name_map)
		selected_instance := prompt.GetUserSelectionFromList("Select EC Instance", instances)

		instance_id := instanceid_and_name_map[selected_instance]

		input := &ssm.StartSessionInput{Target: aws.String(instance_id)}

		ssm_client := ssm.NewFromConfig(cfg)
		sess, err := ssm_client.StartSession(context.Background(), input)
		if err != nil {
			log.Fatalf("%v", err)
		}

		sessJson, _ := json.Marshal(sess)
		inputJson, _ := json.Marshal(input)
		session.ExecCommand(sessJson, inputJson, region)
	},
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
}
