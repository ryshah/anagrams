// Package config provides configuration management for the anagram service.
// It loads settings from a YAML configuration file (config.yaml) and provides
// default values where appropriate.
//
// The configuration is loaded once and cached using a singleton pattern
// to ensure consistency across the application.
package config

import (
	"log/slog"
	"sync"

	"github.com/spf13/viper"
)

// Config represents the complete application configuration structure.
// It uses struct tags for mapping YAML keys to Go struct fields via Viper.
//
// Configuration sections:
//   - Service: General service information (name)
//   - Dictionary: Dictionary file paths for anagram lookup
//   - LRUCache: Cache configuration (capacity)
//   - Server: HTTP server settings (port)
//   - Client: Test client settings (concurrent requests, test words)
//   - Log: Logging configuration (debug mode)
type Config struct {
	// Service contains general service metadata
	Service struct {
		Name string `mapstructure:"name"` // Service name for identification
	} `mapstructure:"service"`

	// Dictionary specifies the word list files to load
	Dictionary struct {
		Files []string `mapstructure:"files"` // Paths to dictionary files
	} `mapstructure:"dictionary"`

	// LRUCache configures the in-memory cache
	LRUCache struct {
		Capacity int `mapstructure:"capacity"` // Maximum number of cached entries
	} `mapstructure:"lru_cache"`

	// Server configures the HTTP server
	Server struct {
		Port string `mapstructure:"port"` // Server port (e.g., ":8080")
	} `mapstructure:"server"`

	// Client configures the test client behavior
	Client struct {
		ConcurrentRequests int      `mapstructure:"concurrent_requests"` // Number of concurrent test requests
		TestWords          []string `mapstructure:"test_words"`          // Words to use in testing
	} `mapstructure:"client"`

	// Log configures logging behavior
	Log struct {
		Debug bool `mapstructure:"debug"` // Enable debug-level logging
	} `mapstructure:"log"`
}

var (
	// cfg is the singleton configuration instance
	cfg Config
	// once ensures configuration is loaded only once
	once sync.Once
	// usingDefaults indicates whether default values were used due to missing/invalid config
	usingDefaults bool
)

// Load reads and returns the application configuration.
// It uses a singleton pattern to ensure the config file is read only once,
// even when called from multiple goroutines.
//
// Configuration loading process:
//  1. Looks for config.yaml in the current directory
//  2. Applies default values for missing settings
//  3. Unmarshals YAML into the Config struct
//  4. Returns a pointer to the singleton config instance and a flag indicating if defaults were used
//
// Default values:
//   - server.port: ":8080"
//   - lru_cache.capacity: 100
//   - dictionary.files: ["data/english.txt"]
//   - service.name: "Anagram Finder"
//
// Returns:
//   - *Config: Pointer to the loaded configuration
//   - bool: true if defaults were used (config file not found or invalid), false if config file was loaded successfully
//
// Example:
//
//	cfg, usedDefaults := config.Load()
//	if usedDefaults {
//	    log.Warn("Using default configuration")
//	}
//	fmt.Println("Server port:", cfg.Server.Port)
func Load() (*Config, bool) {
	// Initialize the config file only once using sync.Once
	// This ensures thread-safe singleton behavior
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		viper.AddConfigPath(".")

		// Set comprehensive defaults
		viper.SetDefault("server.port", ":8080")
		viper.SetDefault("lru_cache.capacity", 100)
		viper.SetDefault("dictionary.files", []string{"data/english.txt"})
		viper.SetDefault("service.name", "Anagram Finder")
		viper.SetDefault("log.debug", false)
		viper.SetDefault("client.concurrent_requests", 10)
		viper.SetDefault("client.test_words", []string{"listen", "silent"})

		err := viper.ReadInConfig()
		if err != nil {
			usingDefaults = true
			slog.Warn("Config file not found, using defaults: " + err.Error())
		}

		err = viper.Unmarshal(&cfg)
		if err != nil {
			usingDefaults = true
			slog.Error("Unable to decode config, using defaults: " + err.Error())
		}
	})
	return &cfg, usingDefaults
}
