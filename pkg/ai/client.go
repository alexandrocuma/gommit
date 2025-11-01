package ai

import (
	"context"
	"fmt"
	"strings"

	"gommit/internal/config"
	"gommit/pkg/ai/providers"
)

// Client implements the AI operations using the configured provider
type Client struct {
	provider providers.Provider
	cfg      *config.AI
}

// NewClient creates a new AI client
func NewClient(cfg *config.Config) (*Client, error) {
	provider, err := NewProvider(&cfg.AI)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI provider: %w", err)
	}

	return &Client{
		provider: provider,
		cfg:      &cfg.AI,
	}, nil
}

// GenerateCommitMessage creates a commit message using the configured AI provider
func (c *Client) GenerateCommitMessage(diff string, data []string, commitCfg *config.Commit) (string, error) {
	prompt := c.buildCommitPrompt(diff, data, commitCfg)

	messages := []providers.Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant that generates clear, concise git commit messages based on code changes.",
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
	message = c.postProcessCommitMessage(message, commitCfg)

	return message, nil
}

// buildCommitPrompt constructs the prompt for commit message generation
func (c *Client) buildCommitPrompt(diff string, data []string, commitCfg *config.Commit) string {
	var prompt strings.Builder

	prompt.WriteString("Generate a git commit message based on the following changes:\n\n")

	// Add context if available
	if len(data) > 0 {
		prompt.WriteString("Context:\n")
		for _, ctx := range data {
			prompt.WriteString("- ")
			prompt.WriteString(ctx)
			prompt.WriteString("\n")
		}
		prompt.WriteString("\n")
	}

	// Add requirements based on configuration
	prompt.WriteString("Requirements:\n")
	if commitCfg.Conventional {
		prompt.WriteString("- Use conventional commit format (type: description)\n")
		prompt.WriteString("- Common types: feat, fix, docs, style, refactor, test, chore\n")
	}

	if commitCfg.Emoji {
		prompt.WriteString("- Include relevant gitmoji at the start\n")
	}

	prompt.WriteString("- Be descriptive but concise\n")
	prompt.WriteString("- Focus on the what and why, not how\n")
	prompt.WriteString("- Maximum 72 characters for subject line\n")

	if commitCfg.Language != "english" {
		prompt.WriteString(fmt.Sprintf("- Write in %s\n", commitCfg.Language))
	}

	prompt.WriteString("\nCode changes (git diff):\n")
	prompt.WriteString("```diff\n")
	prompt.WriteString(diff)
	prompt.WriteString("\n```\n\n")
	prompt.WriteString("Commit message:")

	return prompt.String()
}

// postProcessCommitMessage cleans up the AI-generated commit message
func (c *Client) postProcessCommitMessage(message string, commitCfg *config.Commit) string {
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

// GeneratePRDescription generates a PR description using AI
func (c *Client) GeneratePRDescription(title string, commits []string, diffStats string) (string, error) {
	prompt := c.buildPRPrompt(title, commits, diffStats)

	messages := []providers.Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant that generates comprehensive pull request descriptions.",
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

	return strings.TrimSpace(resp.Content), nil
}

func (c *Client) buildPRPrompt(title string, commits []string, diffStats string) string {
	return fmt.Sprintf(`Generate a comprehensive PR description for: "%s"

Commits in this PR:
%s

Changes:
%s

Please provide:
1. Brief overview of changes
2. Key changes made
3. Testing performed
4. Any breaking changes
5. Additional notes

Format in markdown:`, title, strings.Join(commits, "\n"), diffStats)
}

// GeneratePRDescriptionWithTemplate generates PR description using a template
func (c *Client) GeneratePRDescriptionWithTemplate(title string, commits []string, diff string, diffStats string, template string) (string, error) {
	prompt := c.buildPRTemplatePrompt(title, commits, diff, diffStats, template)

	messages := []providers.Message{
		{
			Role:    "system",
			Content: `
				---
				Role: 
					You are an AI assistant specialized in generating high-quality, consistent, and concise pull request (PR) descriptions. Your primary goal is to create descriptions that are incredibly helpful for reviewers, can be read and understood in less than 5 minutes, and serve as a valuable historical record.
				---
				Core Principles for the PR Description:
					- Consistency: Always follow the provided template structure exactly. Do not invent new sections.
					- Conciseness: Be direct and to the point. Use bullet points and avoid long, dense paragraphs. The total description should be scannable quickly.
					- Helpfulness: Focus on the what, why, and how. Provide context that isn't immediately obvious from the code changes themselves.

				---
				Instructions:
					- You will be given information about the code changes (e.g., file paths, a code diff, or a summary from the developer).
					- Your task is to synthesize this information and populate the following template.
					- Use clear, simple language. Assume the reviewer is skilled but may not have full context on the task.
					- For the "Files Changed" section, group files logically (e.g., "Frontend Components," "API Endpoints," "Database Migrations") rather than listing every file individually.
			`,
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

	return strings.TrimSpace(resp.Content), nil
}

func (c *Client) buildPRTemplatePrompt(title string, commits []string, diff string, diffStats string, template string) string {
	return fmt.Sprintf(`Generate a pull request description using the following template structure.

PR Title: %s

Commits in this PR:
%s

Change Statistics:
%s

Code Changes:
%s

Template to follow (fill in the sections, keep the markdown structure):
%s

Instructions:
- Fill in each section of the template with relevant information
- Keep the same markdown headers and structure
- Be concise but comprehensive
- For changelog, focus on user-facing changes
- For test evidence, mention what was tested and how
- Remove any placeholder comments like <!-- -->`,
		title,
		strings.Join(commits, "\n"),
		diffStats,
		"```diff\n"+diff+"\n```",
		template)
}

func (c *Client) GenerateReview(diff string) (string, error) {
	messages := []providers.Message{
		{
			Role:    "system",
			Content: `
				**Role:** You are a senior engineer conducting a final pre-merge code review. Your goal is to be objective, constructive, and prioritize important issues over nitpicks.

				## Review Framework

				**R - Redundant Code**  
				Unused variables, dead code, duplicated logic, or unnecessary comments that should be removed.

				**I - Improvements & Best Practices**  
				Meaningful optimizations for readability, performance, or adherence to language/framework conventions.

				**S - Security & Logic**  
				Bugs, unhandled edge cases, potential crashes, race conditions, or security vulnerabilities.

				**E - Explanations & Messages**  
				Unclear comments, log messages, commit messages, or documentation that needs improvement.

				**N - Nits & Non-Essentials**  
				Minor formatting, naming preferences, or trivial suggestions that don't affect functionality.

				## Output Format

				1. **Acknowledge:** Briefly state the code's purpose
				2. **Categorize:** List findings under relevant RISEN letters (omit empty categories)
				3. **Justify:** For each point, explain WHY and provide specific suggestions
				4. **Final Verdict:**
					- 🚨 'Changes required before merge' - for any Security/Logic issues
					- 💡 'No critical blockers, but consider suggestions' - for R/I/E/N only
					- ✅ 'No suggestions. Code is ready to merge' - if no issues found

				## Important Guidelines
				- Be concise and practical
				- Distinguish clearly between requirements and suggestions
				- Don't over-complicate or introduce unnecessary friction
				- Focus on what matters for code quality and maintainability
			`,
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