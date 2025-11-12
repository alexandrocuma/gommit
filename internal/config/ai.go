package config

import (
	"fmt"
	"strings"
)

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
	default:
		return fmt.Errorf("unsupported AI provider: %s", c.AI.Provider)
	}

	return nil
}