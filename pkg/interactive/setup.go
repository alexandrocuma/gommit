package interactive

import (
	"fmt"

	"gommit/internal/config"
)

func RunSetup() (*config.Config, error) {
	fmt.Println("🚀 Welcome to Gommit CLI Setup!")
	fmt.Println("Let's configure your assistant...")

	cfg := config.DefaultConfig()

	AI, err := RunAISetup()
	if err != nil {
		return nil, fmt.Errorf("failed to setup ai provider: %w", err)
	}
	cfg.AI = *AI

	err = RunPromptSetup()
	if err != nil {
		return nil, fmt.Errorf("failed to setup prompt templates: %w", err)
	}

	return cfg, nil
}
