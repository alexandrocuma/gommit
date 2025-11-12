package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gommit/internal/config"
	"gommit/pkg/ai/providers"
	"gommit/pkg/utils"
)

// Client implements the AI operations using the configured provider
type Client struct {
	provider providers.Provider
	cfg      *config.AI
	configDir string 
}

// NewClient creates a new AI client
func NewClient(cfg *config.Config) (*Client, error) {
	provider, err := NewProvider(&cfg.AI)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI provider: %w", err)
	}

	// Determine config directory
	configDir, err := utils.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	return &Client{
		provider:  provider,
		cfg:       &cfg.AI,
		configDir: configDir,
	}, nil
}


// loadPromptFile loads a prompt file from the config directory
// Returns empty string if file doesn't exist (use default behavior)
func (c *Client) loadPromptFile(filename string) (string, error) {
	promptPath := filepath.Join(c.configDir, filename)
	
	data, err := os.ReadFile(promptPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // File doesn't exist, use default
		}
		return "", fmt.Errorf("failed to read prompt file %s: %w", promptPath, err)
	}
	
	return string(data), nil
}

// GenerateCommitMessage creates a commit message using the configured AI provider
func (c *Client) GenerateCommitMessage(diff string, data []string) (string, error) {
	promptContent, err := c.loadPromptFile("commit.md")
	if err != nil {
		return "", err
	}

	prompt := c.buildCommitData(diff, data)

	messages := []providers.Message{
		{
			Role:    "system",
			Content: promptContent,
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	req := &providers.ChatRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		Temperature: c.cfg.Temperature,
		MaxTokens:   c.cfg.MaxTokens,
	}

	ctx := context.Background()
	resp, err := c.provider.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("AI completion failed: %w", err)
	}

	// Clean up the response
	message := strings.TrimSpace(resp.Content)
	message = c.postProcessCommitMessage(message)

	return message, nil
}


func (c *Client) buildCommitData(diff string, data []string) string {
	var contextSection string
	if len(data) > 0 {
		items := make([]string, len(data))
		for i, item := range data {
			items[i] = "- " + item
		}
		contextSection = "Context:\n" + strings.Join(items, "\n") + "\n\n"
	}

	return fmt.Sprintf(`%s Diff: %s`, contextSection, "```diff\n"+diff+"\n```")
}

// postProcessCommitMessage cleans up the AI-generated commit message
func (c *Client) postProcessCommitMessage(message string) string {
	// Remove quotes if present
	message = strings.Trim(message, "\"'`")

	// Remove any prefix like "Commit message:"
	if idx := strings.Index(message, ":"); idx != -1 {
		prefix := strings.ToLower(message[:idx])
		if strings.Contains(prefix, "commit") {
			message = strings.TrimSpace(message[idx+1:])
		}
	}

	// Ensure it's a single line (take first line only)
	if lines := strings.Split(message, "\n"); len(lines) > 0 {
		message = strings.TrimSpace(lines[0])
	}

	// Truncate to reasonable length
	if len(message) > 200 {
		message = message[:200] + "..."
	}

	return message
}

// GeneratePRDescriptionWithTemplate generates PR description using a template
func (c *Client) GeneratePRDescriptionWithTemplate(title string, commits []string, diff string, diffStats string, template string) (string, error) {
	prompt, err := c.loadPromptFile("draft.md")
	if err != nil {
		return "", err
	}
	if prompt == "" {
		return "", fmt.Errorf("prompt is missing, check your 'pr description generator' prompt file (draft.md)")
	}

	data := c.buildPRDescriptionData(title, commits, diff, diffStats, template)

	messages := []providers.Message{
		{
			Role: "system",
			Content: prompt,
		},
		{
			Role:    "user",
			Content: data,
		},
	}

	req := &providers.ChatRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		Temperature: c.cfg.Temperature,
		MaxTokens:   c.cfg.MaxTokens,
	}

	ctx := context.Background()
	resp, err := c.provider.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("AI completion failed: %w", err)
	}

	return strings.TrimSpace(resp.Content), nil
}

func (c *Client) buildPRDescriptionData(title string, commits []string, diff string, diffStats string, template string) string {
	return fmt.Sprintf(`PR Title: %s

		Commits in this PR:
		%s

		Change Statistics:
		%s

		Code Changes:
		%s

		Template to follow (fill in the sections, keep the markdown structure):
		%s`,
		title,
		strings.Join(commits, "\n"),
		diffStats,
		"```diff\n"+diff+"\n```",
		template)
}

func (c *Client) GeneratePRReview(diff string) (string, error) {
	prompt, err := c.loadPromptFile("review.md")
	if err != nil {
		return "", err
	}

	if prompt == "" {
		return "", fmt.Errorf("prompt is missing, check your 'pr reviewer' prompt file (review.md)")
	}

	messages := []providers.Message{
		{
			Role: "system",
			Content: prompt,
		},
		{
			Role:    "user",
			Content: diff,
		},
	}

	req := &providers.ChatRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		Temperature: c.cfg.Temperature,
		MaxTokens:   c.cfg.MaxTokens,
	}

	ctx := context.Background()
	resp, err := c.provider.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("AI completion failed: %w", err)
	}

	return strings.TrimSpace(resp.Content), nil
}

