package interactive

import (
	"fmt"

	"github.com/alexandrocuma/gommit/internal/config"

	"github.com/manifoldco/promptui"
)

func RunDirectorySetup() (*config.Directory, error) {
	cfg := config.DefaultDirectoryConfig()
	
	promptsDirPrompt := promptui.Prompt{
		Label:   "Prompts Directory",
		Default: cfg.Prompts,
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("directory path required")
			}
			return nil
		},
	}

	promptsDir, err := promptsDirPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("model input failed: %w", err)
	}

	cfg.Prompts = promptsDir

	templatesDirPrompt := promptui.Prompt{
		Label:   "Templates Directory",
		Default: cfg.Templates,
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("directory path required")
			}
			return nil
		},
	}

	templatesDir, err := templatesDirPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("model input failed: %w", err)
	}
	
	cfg.Templates = templatesDir

	return cfg, nil
}
