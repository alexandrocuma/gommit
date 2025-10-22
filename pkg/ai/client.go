package ai

import (
	"fmt"
	"strings"

	"gommit/internal/config"
)

// Client interface for AI operations
type Client interface {
	GenerateCommitMessage(diff string, context []string, cfg *config.Config) (string, error)
}

// DefaultClient implements the AI Client interface
type DefaultClient struct{}

// NewClient creates a new AI client based on configuration
func NewClient(cfg *config.Config) Client {
	// For now, return the default client
	// In the future, this could return different clients based on cfg.AI.Provider
	return &DefaultClient{}
}

// GenerateCommitMessage creates a commit message using AI
func (c *DefaultClient) GenerateCommitMessage(diff string, context []string, cfg *config.Config) (string, error) {
	// This is a simplified version that creates a basic commit message
	// In a real implementation, you would call the actual AI API here
	
	// Analyze the diff to determine the type of changes
	changeType := analyzeChanges(diff)
	
	// Build commit message based on configuration
	var message strings.Builder
	
	if cfg.Commit.Conventional {
		message.WriteString(getConventionalPrefix(changeType, cfg.Commit.Emoji))
		message.WriteString(": ")
	}
	
	message.WriteString(generateDescriptiveMessage(diff, context))
	
	return message.String(), nil
}

// analyzeChanges determines the type of changes from the diff
func analyzeChanges(diff string) string {
	if strings.Contains(diff, "+++ b/") && strings.Contains(diff, "--- a/") {
		if strings.Count(diff, "+++ b/") > 5 {
			return "feat"
		}
	}
	
	if strings.Contains(diff, "import") || strings.Contains(diff, "package") {
		return "refactor"
	}
	
	if strings.Contains(diff, "fix") || strings.Contains(diff, "bug") || strings.Contains(diff, "error") {
		return "fix"
	}
	
	if strings.Contains(diff, "test") || strings.Contains(diff, "spec") {
		return "test"
	}
	
	return "feat"
}

// getConventionalPrefix returns the conventional commit prefix
func getConventionalPrefix(changeType string, useEmoji bool) string {
	prefixes := map[string]string{
		"feat":     "feat",
		"fix":      "fix", 
		"docs":     "docs",
		"style":    "style",
		"refactor": "refactor",
		"test":     "test",
		"chore":    "chore",
	}
	
	emojis := map[string]string{
		"feat":     "✨",
		"fix":      "🐛",
		"docs":     "📚",
		"style":    "💎",
		"refactor": "🔨",
		"test":     "🧪",
		"chore":    "🔧",
	}
	
	prefix := prefixes[changeType]
	if prefix == "" {
		prefix = "feat"
	}
	
	if useEmoji {
		emoji := emojis[changeType]
		if emoji != "" {
			return fmt.Sprintf("%s %s", emoji, prefix)
		}
	}
	
	return prefix
}

// generateDescriptiveMessage creates a descriptive commit message
func generateDescriptiveMessage(diff string, context []string) string {
	// Simple heuristic-based message generation
	// In real implementation, this would use AI
	
	lines := strings.Split(diff, "\n")
	fileCount := 0
	addedLines := 0
	deletedLines := 0
	
	for _, line := range lines {
		if strings.HasPrefix(line, "+++ b/") {
			fileCount++
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			addedLines++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			deletedLines++
		}
	}
	
	if fileCount == 1 {
		for _, line := range lines {
			if strings.HasPrefix(line, "+++ b/") {
				filename := strings.TrimPrefix(line, "+++ b/")
				return fmt.Sprintf("update %s (%d additions, %d deletions)", filename, addedLines, deletedLines)
			}
		}
	}
	
	return fmt.Sprintf("update %d files (%d additions, %d deletions)", fileCount, addedLines, deletedLines)
}