package interactive

import (
	"fmt"
	"gommit/internal/config"
	"gommit/internal/helpers"

	"github.com/manifoldco/promptui"
)

// RunSetup runs the interactive configuration setup
func RunCommitSetup() (*config.Commit, error) {
	cfg := config.DefaultCommitConfig()
	// Commit Style - Conventional Commits
	convPrompt := promptui.Select{
		Label: "Use Conventional Commits format?",
		Items: []string{"Yes", "No"},
	}

	convIndex, _, err := convPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("conventional commits selection failed: %w", err)
	}
	cfg.Conventional = helpers.IndexToBool(convIndex)

	// Commit Style - Emojis
	emojiPrompt := promptui.Select{
		Label: "Use emojis in commit messages?",
		Items: []string{"Yes", "No"},
	}

	emojiIndex, _, err := emojiPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("emoji selection failed: %w", err)
	}
	cfg.Emoji = helpers.IndexToBool(emojiIndex)

	// Commit Language
	langPrompt := promptui.Select{
		Label: "Commit message language",
		Items: []string{"english", "spanish", "french", "german", "chinese", "japanese"},
	}

	_, lang, err := langPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("language selection failed: %w", err)
	}
	cfg.Language = lang

	return cfg, nil
}
