package cmd

import (
	"context"
	"eclogin/pkg/aws/config"
	"eclogin/pkg/aws/ec2"
	"eclogin/pkg/aws/session"
	"eclogin/pkg/prompt"
	"encoding/json"
	"fmt"
	"log"

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
	Run: func(cmd *cobra.Command, _ []string) {
		region := prompt.GetFlagOrInput(cmd, "region", "Please enter AWS region (default: ap-northeast-1)", "ap-northeast-1")
		profile := prompt.GetFlagOrInput(cmd, "profile", "Please enter AWS profile (optional)", "")

		cfg, err := config.LoadConfig(region, profile)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		ec2Client := aws_ec2.NewFromConfig(cfg)
		instanceNameIDMap := ec2.GetInstanceNameIDMap(ec2Client)
		displayNames := ec2.GetInstanceDisplayNames(instanceNameIDMap)
		if len(displayNames) == 0 {
			log.Fatalf("No EC2 instances found")
		}

		selectedInstance := prompt.GetFlagOrSelect(cmd, "instance-id", "Select EC2 Instance", displayNames)
		instanceID := instanceNameIDMap[selectedInstance]
		printAwsCliEc2Command(instanceID, region)

		sessionInput := &ssm.StartSessionInput{Target: aws.String(instanceID)}
		ssmClient := ssm.NewFromConfig(cfg)

		sessionOutput, err := ssmClient.StartSession(context.Background(), sessionInput)
		if err != nil {
			log.Fatalf("Failed to start SSM session: %v", err)
		}

		sessionData, err := json.Marshal(sessionOutput)
		if err != nil {
			log.Fatalf("Failed to marshal session data: %v", err)
		}

		inputData, err := json.Marshal(sessionInput)
		if err != nil {
			log.Fatalf("Failed to marshal input data: %v", err)
		}

		if err := session.StartSession(sessionData, inputData, region); err != nil {
			log.Fatalf("Failed to start plugin session: %v", err)
		}
	},
}

func printAwsCliEc2Command(instanceID, region string) {
	fmt.Printf(`AWS CLI equivalent command:
aws ssm start-session \
	--target %s \
	--region %s
`,
		instanceID, region)
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
}
