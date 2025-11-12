package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type GitOperations interface {
	IsGitRepository() bool
	GetStagedDiff() (string, error)
	GetCurrentBranch() (string, error)
	GetRecentCommits(count int) ([]string, error)
	Commit(message string) error
	// NEW METHODS FOR PR FEATURE
	GetDiffBetweenBranches(baseBranch, compareBranch string) (string, error)
	GetCommitsBetweenBranches(baseBranch, compareBranch string) ([]string, error)
	GetDiffStatsBetweenBranches(baseBranch, compareBranch string) (string, error)
	BranchExists(branch string) bool
}

type RealGitOperations struct{}

// Helper functions
func (g *RealGitOperations) GetDefaultBaseBranch() string {
	// Try common base branch names
	possibleBranches := []string{"main", "master", "production"}

	for _, branch := range possibleBranches {
		if g.BranchExists(branch) {
			return branch
		}
	}

	// Fallback to main
	return "main"
}

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

	diff := string(output)
	if diff == "" {
		fmt.Println("❌ No staged changes found.")
		fmt.Println("   Please stage your changes first: git add <files>")
		os.Exit(1)
	}

	return diff, nil
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
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}

// NEW: Get diff between two branches
func (g *RealGitOperations) GetDiffBetweenBranches(baseBranch, compareBranch string) (string, error) {
	cmd := exec.Command("git", "diff", fmt.Sprintf("%s..%s", baseBranch, compareBranch))

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff between branches: %w", err)
	}

	diff := string(output)
	if diff == "" {
		fmt.Println("❌ No staged changes found.")
		fmt.Println("   Please stage your changes first: git add <files>")
		os.Exit(1)
	}

	return diff, nil
}

// NEW: Get commit history between branches
func (g *RealGitOperations) GetCommitsBetweenBranches(baseBranch, compareBranch string) ([]string, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("%s..%s", baseBranch, compareBranch), "--oneline", "--no-decorate")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commits between branches: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

// NEW: Get diff statistics between branches
func (g *RealGitOperations) GetDiffStatsBetweenBranches(baseBranch, compareBranch string) (string, error) {
	cmd := exec.Command("git", "diff", "--stat", fmt.Sprintf("%s..%s", baseBranch, compareBranch))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff stats between branches: %w", err)
	}
	return string(output), nil
}

// NEW: Check if branch exists
func (g *RealGitOperations) BranchExists(branch string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", fmt.Sprintf("refs/heads/%s", branch))
	return cmd.Run() == nil
}
