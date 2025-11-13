package interactive

import (
	"fmt"
	"gommit/internal/config"
	"gommit/internal/helpers"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

func RunAISetup() (*config.AI, error) {
	cfg := config.DefaultAIConfig()
	providerPrompt := promptui.Select{
		Label: "Select AI Provider",
		Items: []string{"openai", "anthropic", "deepseek"},
	}

	_, provider, err := providerPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("provider selection failed: %w", err)
	}

	cfg.Provider = provider

	apiKeyPrompt := promptui.Prompt{
		Label: fmt.Sprintf("Enter your %s API Key", strings.ToUpper(provider)),
		Mask:  '*',
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("API key is required")
			}
			return nil
		},
	}

	apiKey, err := apiKeyPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("API key input failed: %w", err)
	}
	cfg.APIKey = apiKey

	modelPrompt := promptui.Prompt{
		Label:   "AI Model",
		Default: cfg.Model,
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("model is required")
			}
			return nil
		},
	}

	model, err := modelPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("model input failed: %w", err)
	}
	cfg.Model = model

	// Temperature
	tempPrompt := promptui.Prompt{
		Label:    "Temperature (0.0 - 1.0)",
		Default:  fmt.Sprintf("%.1f", cfg.Temperature),
		Validate: helpers.ValidateTemperature,
	}

	tempStr, err := tempPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("temperature input failed: %w", err)
	}
	temp, _ := strconv.ParseFloat(tempStr, 64)
	cfg.Temperature = temp

	// Temperature
	tokensPrompt := promptui.Prompt{
		Label:    "Maximum token count usage (1 - 4096+)",
		Default:  fmt.Sprintf("%d", cfg.MaxTokens),
	}

	tokenStr, err := tokensPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("token count input failed: %w", err)
	}
	// Parse as int64, then convert to int
	tokens, err := strconv.ParseInt(tokenStr, 10, 0) // Use base 10 for decimal input
	if err != nil {
			return nil, fmt.Errorf("invalid token count value: %w", err)
	}

	cfg.MaxTokens = int(tokens) // Convert int64 to int
	
	return cfg, nil
}
