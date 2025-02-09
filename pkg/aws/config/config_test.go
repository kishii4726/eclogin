package config_test

import (
	"eclogin/pkg/aws/config"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	region := "ap-northeast-1"
	profile := "default"

	cfg, _ := config.LoadConfig(region, profile)

	if cfg.Region != region {
		t.Errorf("expected region %s, but got %s", region, cfg.Region)
	}
}
