// Package config provides unit tests for configuration loading.
// Tests verify that default values are properly applied when config file is missing.
package config

import (
	"testing"

	"github.com/spf13/viper"
)

// TestLoad_Defaults verifies that the configuration loader applies default values
// when no config.yaml file is present.
//
// Test scenario:
//  1. Resets Viper to clear any previously loaded configuration
//  2. Calls Load() which should apply defaults
//  3. Verifies that the server port has a default value set
//
// This ensures the application can start with sensible defaults even without
// a configuration file.
func TestLoad_Defaults(t *testing.T) {

	// Reset viper to ensure clean state for testing
	viper.Reset()

	cfg, usedDefaults := Load()
	if !usedDefaults {
		t.Error("expected defaults to be used when config file not found")
	}

	// Verify that default server port is set
	if cfg.Server.Port == "" {
		t.Fatalf("expected server port to be set")
	}
}
