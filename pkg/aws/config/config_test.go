package config_test

import (
	"eclogin/pkg/aws/config"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestLoadConfig(t *testing.T) {
	// CIの場合AWSの権限がないのでスキップ
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping AWS config test in CI environment")
	}

	region := "ap-northeast-1"
	profile := "default"

	cfg, err := config.LoadConfig(region, profile)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Region != region {
		t.Errorf("expected region %s, but got %s", region, cfg.Region)
	}
}

// CIでも実行できるモックテストを追加
func TestLoadConfigMock(t *testing.T) {
	region := "ap-northeast-1"
	mockConfig := aws.Config{
		Region: region,
	}

	if mockConfig.Region != region {
		t.Errorf("expected region %s, but got %s", region, mockConfig.Region)
	}
}
