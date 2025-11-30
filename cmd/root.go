package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	appVersion = "0.0.20"
	appName    = "eclogin"
)

var rootCmd = &cobra.Command{
	Use:     appName,
	Version: appVersion,
	Short:   "CLI tool for logging into AWS EC2/ECS/Local docker containers",
	Long:    `A command-line interface tool that helps you connect to AWS EC2 instances and ECS containers.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Toggle feature flag")

	// EC2 command flags
	ec2Cmd.Flags().StringP("region", "r", "", "AWS region name")
	ec2Cmd.Flags().StringP("profile", "p", "", "AWS profile name")
	ec2Cmd.Flags().StringP("instance-id", "i", "", "EC2 instance ID")

	// ECS command flags
	ecsCmd.Flags().StringP("region", "r", "", "AWS region name")
	ecsCmd.Flags().StringP("profile", "p", "", "AWS profile name")
	ecsCmd.Flags().StringP("cluster", "c", "", "ECS cluster name")
	ecsCmd.Flags().StringP("service", "s", "", "ECS service name")
	ecsCmd.Flags().StringP("task-id", "t", "", "ECS task ID")
	ecsCmd.Flags().StringP("container", "C", "", "ECS container name")
	ecsCmd.Flags().StringP("shell", "S", "", "Shell to use for the session")
}
