package main

import (
	"testing"
)

func TestLoadEnvConfig(t *testing.T) {
	config, err := loadEnvConfig()
	if err != nil {
		t.Error("failed to load the env file.")
	}
	if config.CSRF.Key == "" {
		t.Error("csrf key cannot be blank")
	}
}
