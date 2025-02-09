package prompt

import (
	"testing"

	"github.com/spf13/cobra"
)

type MockPrompter struct {
	inputResult  string
	selectResult string
}

func (m *MockPrompter) Input(label string, defaultValue string) string {
	return m.inputResult
}

func (m *MockPrompter) Select(label string, options []string) string {
	return m.selectResult
}

func TestGetFlagOrInput(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("test-flag", "", "test flag")

	err := cmd.Flags().Set("test-flag", "flag-value")
	if err != nil {
		t.Fatal(err)
	}
	result := GetFlagOrInput(cmd, "test-flag", "Enter value:", "default", &MockPrompter{})
	if result != "flag-value" {
		t.Errorf("Expected flag-value, got %s", result)
	}
}

func TestGetFlagOrSelect(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("test-flag", "", "test flag")

	err := cmd.Flags().Set("test-flag", "option1")
	if err != nil {
		t.Fatal(err)
	}
	options := []string{"option1", "option2", "option3"}
	result := GetFlagOrSelect(cmd, "test-flag", "Select option:", options, &MockPrompter{})
	if result != "option1" {
		t.Errorf("Expected option1, got %s", result)
	}
}

func TestPromptInput(t *testing.T) {
	mock := &MockPrompter{inputResult: "test input"}
	result := mock.Input("Test Label", "default")
	if result != "test input" {
		t.Errorf("Expected 'test input', got %s", result)
	}
}

func TestPromptSelect(t *testing.T) {
	mock := &MockPrompter{selectResult: "option1"}
	options := []string{"option1", "option2", "option3"}
	result := mock.Select("Test Label", options)
	if result != "option1" {
		t.Errorf("Expected 'option1', got %s", result)
	}
}
