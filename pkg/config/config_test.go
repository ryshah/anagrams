package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestLoad_Defaults(t *testing.T) {

	// reset viper
	viper.Reset()

	cfg := Load()

	if cfg.Server.Port == "" {
		t.Fatalf("expected server port to be set")
	}
}
