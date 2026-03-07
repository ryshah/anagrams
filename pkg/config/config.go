// Configuration for anagram service
package config

import (
	"log/slog"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Service struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"service"`

	Dictionary struct {
		Files []string `mapstructure:"files"`
	} `mapstructure:"dictionary"`

	LRUCache struct {
		Capacity int `mapstructure:"capacity"`
	} `mapstructure:"lru_cache"`

	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Client struct {
		ConcurrentRequests int      `mapstructure:"concurrent_requests"`
		TestWords          []string `mapstructure:"test_words"`
	} `mapstructure:"client"`

	Log struct {
		Debug bool `mapstructure:"debug"`
	} `mapstructure:"log"`
}

var (
	cfg  Config
	once sync.Once
)

func Load() *Config {
	// initializes the config file only once
	// provides defaults where appropriate
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		viper.AddConfigPath(".")

		// defaults
		viper.SetDefault("server.port", 8080)

		err := viper.ReadInConfig()
		if err != nil {
			slog.Warn("Config file not found, using env/defaults")
		}

		err = viper.Unmarshal(&cfg)
		if err != nil {
			slog.Error("Unable to decode config: " + err.Error())
		}
	})
	return &cfg
}
