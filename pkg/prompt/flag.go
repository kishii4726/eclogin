package prompt

import "github.com/spf13/cobra"

// HasRequiredFlags checks if all required flags are provided
func HasRequiredFlags(cmd *cobra.Command, flags []string) bool {
	for _, flag := range flags {
		if !cmd.Flags().Changed(flag) {
			return false
		}
	}
	return true
}
