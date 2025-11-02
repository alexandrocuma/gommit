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

	Commit, err := RunCommitSetup()
	if err != nil {
		return nil, fmt.Errorf("failed to setup commit format: %w", err)
	}
	cfg.Commit = *Commit

	return cfg, nil
}
