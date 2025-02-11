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
	Run: runEC2command,
}

func runEC2command(cmd *cobra.Command, _ []string) {
	requiredFlags := []string{"instance-id", "region"}
	prompter := prompt.NewUIPrompter()
	region := prompt.GetFlagOrInput(cmd, "region", "Please enter AWS region (default: ap-northeast-1)", "ap-northeast-1", prompter)

	var profile string
	if prompt.HasRequiredFlags(cmd, requiredFlags) {
		profile = cmd.Flag("profile").Value.String()
	} else {
		profile = prompt.GetFlagOrInput(cmd, "profile", "Please enter AWS profile (optional)", "", prompter)
	}

	cfg, err := config.LoadConfig(region, profile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ec2Client := aws_ec2.NewFromConfig(cfg)

	var instanceID string
	if prompt.HasRequiredFlags(cmd, requiredFlags) {
		instanceID = cmd.Flag("instance-id").Value.String()
	} else {
		instanceNameIDMap := ec2.GetInstanceNameIDMap(ec2Client)
		displayNames := ec2.GetInstanceDisplayNames(instanceNameIDMap)
		if len(displayNames) == 0 {
			log.Fatalf("No EC2 instances found")
		}

		selectedInstance := prompt.GetFlagOrSelect(cmd, "instance-id", "Select EC2 Instance", displayNames, prompter)
		instanceID = instanceNameIDMap[selectedInstance]
		printEcloginEc2WithOptionCommand(cmd, instanceID, region, profile)
	}
	printAwsCliEc2Command(cmd, instanceID, region, profile)

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
}

func printEcloginEc2WithOptionCommand(cmd *cobra.Command, instanceID string, region string, profile string) {
	if !cmd.Flags().Changed("profile") {
		fmt.Printf(`eclogin equivalent command:
eclogin ec2 --instance-id %s --region %s

`,
			instanceID, region)
	} else {
		fmt.Printf(`eclogin equivalent command:
eclogin ec2 --instance-id %s --region %s --profile %s

`,
			instanceID, region, profile)
	}
}

func printAwsCliEc2Command(cmd *cobra.Command, instanceID string, region string, profile string) {
	if !cmd.Flags().Changed("profile") {
		fmt.Printf(`If you are using awscli, please copy the following:
aws ssm start-session \
	--target %s \
	--region %s

`,
			instanceID, region)
	} else {
		fmt.Printf(`If you are using awscli, please copy the following:
aws ssm start-session \
	--target %s \
	--region %s \
	--profile %s

`,
			instanceID, region, profile)
	}
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
}
