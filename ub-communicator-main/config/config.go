// Package config initializes Viper configuration from config.yaml and environment
// variables. Environment variables use the prefix UBCOMMUNICATOR_ and replace
// dots with underscores (e.g., rabbitmq.dsn → UBCOMMUNICATOR_RABBITMQ_DSN).
//
// Configuration precedence: environment variables > config.yaml defaults.
// The config path can be overridden via CONFIG_PATH environment variable.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const FileName = "config"

func SetConfigs() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(FileName)

	// Allow overriding config path via environment variable
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config from %s: %w", configPath, err)
	}

	v.AutomaticEnv()
	v.SetEnvPrefix("ubcommunicator")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v, nil
}
