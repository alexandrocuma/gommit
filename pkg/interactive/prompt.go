package interactive

import (
	"fmt"
	"gommit/internal/config"

	"github.com/manifoldco/promptui"
)

// RunSetup runs the interactive configuration setup
func RunPromptSetup() (*config.Prompt, error) {
	cfg := config.DefaultPromptConfig()
	// Prompt files - Commit Messages, PR Reviews and PR Descriptions
	
	commitPrompt := promptui.Prompt{
		Label:   "Commit - Prompt File",
		Default: cfg.Commit,
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("model is required")
			}
			return nil
		},
	}

	commit, err := commitPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("model input failed: %w", err)
	}
	cfg.Commit = commit
	
	return cfg, nil
}
