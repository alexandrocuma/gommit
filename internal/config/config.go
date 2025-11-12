package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	AI     AI     `yaml:"ai" mapstructure:"ai"`
	Commit Commit `yaml:"commit" mapstructure:"commit"`
}

func DefaultConfig() *Config {
	return &Config{
		AI:     *DefaultAIConfig(),
		Commit: *DefaultCommitConfig(),
	}
}

// getDefaultConfigPath returns the preferred default config path: ~/.gommit/.gommit.config.yaml
func getDefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	return filepath.Join(home, ".gommit", ".gommit.config.yaml"), nil
}

// ConfigExists checks if a configuration file exists in standard locations:
// 1. ./gommit/.gommit.config.yaml (relative to current directory)
// 2. ~/.gommit/.gommit.config.yaml (in user's home directory)
func ConfigExists() bool {
	// Check current working directory subdirectory
	cwd, err := os.Getwd()
	if err == nil {
		currentDirPath := filepath.Join(cwd, "gommit", ".gommit.config.yaml")
		_, err := os.Stat(currentDirPath)
		if err == nil {
			return true
		}
	}

	// Check home directory
	homePath, err := getDefaultConfigPath()
	if err == nil {
		_, err := os.Stat(homePath)
		if err == nil {
			return true
		}
	}

	return false
}

// LoadConfig loads configuration from file with the following precedence:
// 1. ./gommit/.gommit.config.yaml
// 2. ~/.gommit/.gommit.config.yaml
// If no config file is found, returns default configuration
func LoadConfig() (*Config, error) {
	// Configure Viper to look for .gommit.config.yaml files
	viper.SetConfigName(".gommit.config")
	viper.SetConfigType("yaml")

	// Add search paths in order of precedence
	viper.AddConfigPath("./gommit") // Will look for ./gommit/.gommit.config.yaml
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(home, ".gommit")) // Will look for ~/.gommit/.gommit.config.yaml
	}

	// Set default values
	viper.SetDefault("ai", DefaultAIConfig())
	viper.SetDefault("commit", DefaultCommitConfig())

	// Attempt to read config file
	err = viper.ReadInConfig();
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			// Config file not found, return defaults
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal config into struct
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves configuration to file.
// If a config file was previously loaded, it saves to the same location.
// Otherwise, it saves to ~/.gommit/.gommit.config.yaml
func SaveConfig(cfg *Config) error {
	viper.Set("ai", cfg.AI)
	viper.Set("commit", cfg.Commit)

	// Determine where to save
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		// No config was loaded, use default home directory location
		var err error
		configPath, err = getDefaultConfigPath()
		if err != nil {
			return err
		}
	}

	// Ensure target directory exists
	configDir := filepath.Dir(configPath)
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating config directory %q: %w", configDir, err)
	}

	// Write config file
	err = viper.WriteConfigAs(configPath)
	if err != nil {
		return fmt.Errorf("error writing config file %q: %w", configPath, err)
	}

	return nil
}