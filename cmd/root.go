package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eclogin",
	Short: "Checks if the security group and WEF and ALB contain the specified IP address",
	Long: `Checks if the security group and WAF and ALB contain the specified IP address.
You can specify IP addresses as arguments or read a csv file.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	ec2Cmd.Flags().StringP("region", "r", "", "aws region name")
	ec2Cmd.Flags().StringP("profile", "p", "", "aws profile name")
	ec2Cmd.Flags().StringP("instance-id", "i", "", "ec2 instance id")

	ecsCmd.Flags().StringP("region", "r", "", "aws region name")
	ecsCmd.Flags().StringP("profile", "p", "", "aws profile name")
	ecsCmd.Flags().StringP("cluster", "c", "", "ecs cluster name")
	ecsCmd.Flags().StringP("service", "s", "", "ecs service name")
	ecsCmd.Flags().StringP("task-id", "t", "", "ecs task id")
}
