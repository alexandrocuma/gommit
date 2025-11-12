package config

import (
	"fmt"
	"os"
)

func (c *Config) ValidateAIConfig() {
	if c.AI.APIKey == "" {
		fmt.Println("❌ AI API key not configured.")
		fmt.Println("Please run 'gommit init' to set up your configuration.")
		os.Exit(1)
	}
}
