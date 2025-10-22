package ai

import (
	"fmt"

	"gommit/internal/config"
	"gommit/pkg/ai/providers"
)

// NewProvider creates a new AI provider based on configuration
func NewProvider(cfg *config.AI) (providers.Provider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key is required for provider: %s", cfg.Provider)
	}

	switch cfg.Provider {
	case "openai":
		return providers.NewOpenAIProvider(cfg.APIKey), nil
	case "anthropic":
		return providers.NewAnthropicProvider(cfg.APIKey), nil
	case "deepseek":
		return providers.NewDeepSeekProvider(cfg.APIKey), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.Provider)
	}
}
