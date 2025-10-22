package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type GitOperations interface {
	IsGitRepository() bool
	GetStagedDiff() (string, error)
	GetCurrentBranch() (string, error)
	GetRecentCommits(count int) ([]string, error)
	Commit(message string) error
}

type RealGitOperations struct{}

func (g *RealGitOperations) IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func (g *RealGitOperations) GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}
	return string(output), nil
}

func (g *RealGitOperations) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (g *RealGitOperations) GetRecentCommits(count int) ([]string, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", count), "--oneline", "--no-decorate")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent commits: %w", err)
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

func (g *RealGitOperations) Commit(message string) error {
	// Use -m to avoid editor for AI-generated commits
	cmd := exec.Command("git", "commit", "-m", message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}