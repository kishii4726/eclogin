package prompt

import (
	"log"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func GetUserInputFrom(l string, d string) string {
	prompt := promptui.Prompt{
		Label:   l,
		Default: d,
	}
	v, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return v
}

func GetUserSelectionFromList(l string, i []string) string {
	prompt := promptui.Select{
		Label: l,
		Items: i,
	}
	_, v, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return v
}

func GetFlag(cmd *cobra.Command, flag_name string, description_text string, default_value string) string {
	flag_value, err := cmd.Flags().GetString(flag_name)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if flag_value == "" {
		flag_value = GetUserInputFrom(description_text, default_value)
	}

	return flag_value
}

func GetFlagOrPrompt(cmd *cobra.Command, flag_name string, prompt_message string, get_list_func func() []string) string {
	flag_value, err := cmd.Flags().GetString(flag_name)
	if err != nil {
		log.Fatalf("Get argument --%s failed %v\n", flag_name, err)
	}
	if flag_value != "" {
		return flag_value
	}
	return GetUserSelectionFromList(prompt_message, get_list_func())
}
