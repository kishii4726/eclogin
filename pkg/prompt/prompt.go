package prompt

import (
	"log"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Prompter interface {
	Input(label string, defaultValue string) string
	Select(label string, options []string) string
}

type UIPrompter struct{}

func GetFlagOrInput(cmd *cobra.Command, flagName string, promptMsg string, defaultValue string, prompter Prompter) string {
	flagValue, err := cmd.Flags().GetString(flagName)
	if err != nil {
		log.Fatalf("Failed to get flag '%s': %v", flagName, err)
	}

	if flagValue == "" {
		flagValue = prompter.Input(promptMsg, defaultValue)
	}
	return flagValue
}

func GetFlagOrSelect(cmd *cobra.Command, flagName string, promptMsg string, options []string, prompter Prompter) string {
	flagValue, err := cmd.Flags().GetString(flagName)
	if err != nil {
		log.Fatalf("Failed to get flag '%s': %v", flagName, err)
	}
	if flagValue != "" {
		return flagValue
	}
	return prompter.Select(promptMsg, options)
}

func NewUIPrompter() Prompter {
	return &UIPrompter{}
}

func (p *UIPrompter) Input(label string, defaultValue string) string {
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

func (p *UIPrompter) Select(label string, options []string) string {
	prompt := promptui.Select{
		Label:             label,
		Items:             options,
		StartInSearchMode: true,
		Searcher: func(input string, index int) bool {
			option := options[index]
			// Filter options by checking if the input is contained in the option (case-insensitive)
			return strings.Contains(strings.ToLower(option), strings.ToLower(input))
		},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Failed to get user selection: %v\n", err)
	}
	return result
}
