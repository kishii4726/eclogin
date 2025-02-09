package ec2

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Test_getInstanceName(t *testing.T) {
	tests := []struct {
		name     string
		tags     []types.Tag
		expected string
	}{
		{
			name: "正常系：Nameタグあり",
			tags: []types.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("test-instance"),
				},
			},
			expected: "test-instance",
		},
		{
			name: "正常系：Nameタグなし",
			tags: []types.Tag{
				{
					Key:   aws.String("Environment"),
					Value: aws.String("test"),
				},
			},
			expected: "No Name Tag",
		},
		{
			name:     "正常系：タグなし",
			tags:     []types.Tag{},
			expected: "No Name Tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getInstanceName(tt.tags)
			if result != tt.expected {
				t.Errorf("getInstanceName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetInstanceDisplayNames(t *testing.T) {
	instanceMap := map[string]string{
		"test1(i-123)": "i-123",
		"test2(i-456)": "i-456",
	}

	result := GetInstanceDisplayNames(instanceMap)
	if len(result) != 2 {
		t.Errorf("Expected length 2, got %d", len(result))
	}

	expectedNames := map[string]bool{
		"test1(i-123)": false,
		"test2(i-456)": false,
	}

	for _, name := range result {
		if _, exists := expectedNames[name]; !exists {
			t.Errorf("Unexpected name in results: %s", name)
		}
		expectedNames[name] = true
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected name not found in results: %s", name)
		}
	}
}
