package prompt

import (
	"log"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func promptInput(label string, defaultValue string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Failed to get user input: %v\n", err)
	}
	return result
}

func PromptSelect(label string, options []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: options,
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Failed to get user selection: %v\n", err)
	}
	return result
}

func GetFlagOrInput(cmd *cobra.Command, flagName string, promptMsg string, defaultValue string) string {
	flagValue, err := cmd.Flags().GetString(flagName)
	if err != nil {
		log.Fatalf("Failed to get flag '%s': %v", flagName, err)
	}

	if flagValue == "" {
		flagValue = promptInput(promptMsg, defaultValue)
	}
	return flagValue
}

func GetFlagOrSelect(cmd *cobra.Command, flagName string, promptMsg string, options []string) string {
	flagValue, err := cmd.Flags().GetString(flagName)
	if err != nil {
		log.Fatalf("Failed to get flag '%s': %v", flagName, err)
	}
	if flagValue != "" {
		return flagValue
	}
	return PromptSelect(promptMsg, options)
}
