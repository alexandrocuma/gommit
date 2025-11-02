package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AI     AI     `yaml:"ai" mapstructure:"ai"`
	Commit Commit `yaml:"commit" mapstructure:"commit"`
}

func DefaultConfig() *Config {
	cfg := &Config{}

	cfg.AI = *DefaultAIConfig()
	cfg.Commit = *DefaultCommitConfig()

	return cfg
}

// ConfigExists checks if a configuration file exists
func ConfigExists() bool {
	configPath := viper.ConfigFileUsed()
	if configPath != "" {
		_, err := os.Stat(configPath)
		return err == nil
	}

	// Check common locations
	locations := []string{
		".gommit.config.yaml",
		filepath.Join(os.Getenv("HOME"), ".gommit.config.yaml"),
	}

	for _, loc := range locations {
		_, err := os.Stat(loc);
		if err == nil {
			return true
		}
	}

	return false
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	// Set up Viper
	viper.SetConfigName(".gommit.config")
	viper.SetConfigType("yaml")

	// Search paths: current directory, home directory
	viper.AddConfigPath(".")
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(home)
	}

	// Set defaults
	viper.SetDefault("ai.provider", "openai")
	viper.SetDefault("ai.model", "gpt-4")
	viper.SetDefault("ai.temperature", 0.7)
	viper.SetDefault("ai.max_tokens", 500)
	viper.SetDefault("commit.conventional", true)
	viper.SetDefault("commit.emoji", true)
	viper.SetDefault("commit.language", "english")
	viper.SetDefault("pr.template", "default")
	viper.SetDefault("pr.auto_assign", false)
	viper.SetDefault("pr.include_tests", true)

	// Try to read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	viper.Set("ai.provider", cfg.AI.Provider)
	viper.Set("ai.api_key", cfg.AI.APIKey)
	viper.Set("ai.model", cfg.AI.Model)
	viper.Set("ai.temperature", cfg.AI.Temperature)
	viper.Set("ai.max_tokens", cfg.AI.MaxTokens)
	viper.Set("commit.conventional", cfg.Commit.Conventional)
	viper.Set("commit.emoji", cfg.Commit.Emoji)
	viper.Set("commit.language", cfg.Commit.Language)

	// Determine config file path
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		// Use home directory as default location
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}
		configPath = filepath.Join(home, ".gommit.config.yaml")
	}

	// Ensure the directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Write config file
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

type AI struct {
	Provider    string  `yaml:"provider" mapstructure:"provider"`
	APIKey      string  `yaml:"api_key" mapstructure:"api_key"`
	Model       string  `yaml:"model" mapstructure:"model"`
	Temperature float64 `yaml:"temperature" mapstructure:"temperature"`
	MaxTokens   int     `yaml:"max_tokens" mapstructure:"max_tokens"`
}

func DefaultAIConfig() *AI {
	cfg := &AI{}

	// AI defaults
	cfg.Provider = "openai"
	cfg.Model = "gpt-4"
	cfg.Temperature = 0.7
	cfg.MaxTokens = 500

	return cfg
}

// Add this method to the Config struct
func (c *Config) Validate() error {
	// Validate AI configuration
	if c.AI.APIKey == "" {
		return fmt.Errorf("AI API key is required")
	}

	// Add provider-specific validation if needed
	switch c.AI.Provider {
	case "openai":
		if !strings.HasPrefix(c.AI.APIKey, "sk-") {
			return fmt.Errorf("invalid OpenAI API key format")
		}
	case "anthropic":
		if !strings.HasPrefix(c.AI.APIKey, "sk-ant-") {
			return fmt.Errorf("invalid Anthropic API key format")
		}
	case "deepseek":
		// DeepSeek keys don't have a specific format
	case "azure-openai":
		// Azure keys are typically base64 encoded
	default:
		return fmt.Errorf("unsupported AI provider: %s", c.AI.Provider)
	}

	return nil
}

type Commit struct {
	Conventional bool   `yaml:"conventional" mapstructure:"conventional"`
	Emoji        bool   `yaml:"emoji" mapstructure:"emoji"`
	Language     string `yaml:"language" mapstructure:"language"`
}

func DefaultCommitConfig() *Commit {
	cfg := &Commit{}

	// AI defaults
	cfg.Conventional = true
	cfg.Emoji = false
	cfg.Language = "english"

	return cfg
}
