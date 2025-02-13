package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
)

type MockPrompter struct {
	mock.Mock
}

func (m *MockPrompter) Input(message string, defaultValue string) string {
	args := m.Called(message, defaultValue)
	return args.String(0)
}

func (m *MockPrompter) Select(message string, options []string) string {
	args := m.Called(message, options)
	return args.String(0)
}

func TestPrintAwsCliEc2Command(t *testing.T) {
	tests := []struct {
		cmd        *cobra.Command
		name       string
		instanceID string
		region     string
		profile    string
		expected   string
	}{
		{
			cmd:        &cobra.Command{},
			name:       "prints correct AWS CLI command",
			instanceID: "i-1234567890abcdef0",
			region:     "ap-northeast-1",
			profile:    "",
			expected: `If you are using awscli, please copy the following:
aws ssm start-session \
	--target i-1234567890abcdef0 \
	--region ap-northeast-1

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput(func() {
				printAwsCliEc2Command(tt.cmd, tt.instanceID, tt.region, tt.profile)
			})
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}
