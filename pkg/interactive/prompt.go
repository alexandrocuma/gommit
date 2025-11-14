package interactive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexandrocuma/gommit/pkg/utils"

	"github.com/manifoldco/promptui"
)

// Default prompt content embedded in code
const (
	defaultCommitPrompt = `# Commit Message Generator Prompt
You are an expert Git commit message generator. Analyze the provided code diff and generate a concise, clear commit message following these rules:
- Use present tense ("Add feature" not "Added feature")
- Keep the first line under 50 characters
- Provide a detailed description if the change is complex
- Reference any related issues or tickets
`

	defaultDraftPrompt = `# PR Description Generator Prompt
You are a skilled technical writer crafting pull request descriptions. Based on the code changes provided, create a comprehensive PR description that includes:
- Clear summary of changes
- Motivation and context
- Breaking changes (if any)
- Testing instructions
- Screenshots for UI changes (mention if applicable)
`

	defaultReviewPrompt = `# Code Review Prompt
You are an experienced senior software engineer conducting a thorough code review. Analyze the code changes and provide:
- Potential bugs or errors
- Security vulnerabilities
- Performance issues
- Code style improvements
- Best practice recommendations
- Specific line comments where relevant
Be constructive and educational in your feedback.
`
)

// RunPromptSetup runs the interactive configuration setup for prompt files
func RunPromptSetup() error {
	// Get the configuration directory
	cfgDir, err := utils.GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Ensure the directory exists
	err = os.MkdirAll(cfgDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory %q: %w", cfgDir, err)
	}

	// Define prompts to configure
	prompts := []struct {
		filename     string
		defaultValue string
		label        string
	}{
		{"commit.md", defaultCommitPrompt, "💬 Configure 'commit generator' prompt"},
		{"draft.md", defaultDraftPrompt, "💬 Configure 'PR description generator' prompt"},
		{"review.md", defaultReviewPrompt, "💬 Configure 'PR reviewer' prompt"},
	}

	// Setup each prompt
	for _, p := range prompts {
		err := setupPrompt(cfgDir, p.filename, p.defaultValue, p.label)
		if err != nil {
			return fmt.Errorf("setup failed for %s: %w", p.filename, err)
		}
	}

	return nil
}

// setupPrompt handles the interactive setup for a single prompt file
func setupPrompt(cfgDir, filename, defaultContent, label string) error {
	// Select default or custom
	promptTypeSelect := promptui.Select{
		Label: label,
		Items: []string{"default", "custom"},
	}

	_, choice, err := promptTypeSelect.Run()
	if err != nil {
		return fmt.Errorf("selection failed: %w", err)
	}

	var content string
	if choice == "custom" {
		// Get custom prompt from user
		customPrompt := promptui.Prompt{
			Label:    "Enter your custom prompt",
			Validate: validateNonEmpty,
		}

		content, err = customPrompt.Run()
		if err != nil {
			return fmt.Errorf("custom input failed: %w", err)
		}
	} else {
		content = defaultContent
	}

	// Save to file
	filePath := filepath.Join(cfgDir, filename)
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %q: %w", filePath, err)
	}

	fmt.Printf("✅ Saved %s\n", filepath.Join(filepath.Base(cfgDir), filename))
	return nil
}

// validateNonEmpty ensures input is not empty
func validateNonEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("prompt cannot be empty")
	}
	return nil
}