/*
Copyright © 2025 Alexandro Cu alexandro.cuma@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexandrocuma/gommit/internal/config"
	"github.com/alexandrocuma/gommit/internal/git"
	"github.com/alexandrocuma/gommit/pkg/ai"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	skipConfirm bool
	verbose     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gommit",
	Short: "Generate commit messages for staged changes",
	Long: `Automatically generate meaningful commit messages using AI based on your staged changes.

		This command analyzes your git diff, understands the context of your changes,
		and creates a descriptive commit message that follows conventional commit standards.

		Features:
		• Analyzes staged changes and git context
		• Supports multiple AI providers (OpenAI, Anthropic, etc.)
		• Follows your preferred commit style (conventional, semantic, etc.)
		• Interactive confirmation before committing
		• Configurable base branch comparison

		Examples:
			git add . && gommit       # Commit all staged changes
			gommit --verbose          # Show detailed process
			gommit --no-confirm       # Skip confirmation prompt
			gommit --base main        # Compare against main branch`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("❌ Failed to load configuration: %v", err)
		}

		cfg.ValidateAIConfig()

		if verbose {
			fmt.Printf("🤖 Using AI provider: %s\n", cfg.AI.Provider)
		}

		// Initialize git operations
		gitOps := &git.RealGitOperations{}

		// Check if we're in a git repository
		if !gitOps.IsGitRepository() {
			fmt.Println("❌ Not a git repository")
			os.Exit(1)
		}

		// Get current branch
		currentBranch, err := gitOps.GetCurrentBranch()
		if err != nil {
			log.Fatalf("❌ Failed to get current branch: %v", err)
		}

		// Set default base branch if not provided
		if baseBranch == "" {
			baseBranch = gitOps.GetDefaultBaseBranch()
		}

		fmt.Printf("📊 Comparing changes from '%s' to '%s'...\n", currentBranch, baseBranch)

		if verbose {
			fmt.Println("📊 Analyzing staged changes...")
		}

		diff, err := gitOps.GetStagedDiff()
		if err != nil {
			log.Fatalf("❌ Error getting git diff: %v", err)
		}

		// Get context for better commit messages
		var context []string
		branch, err := gitOps.GetCurrentBranch()
		if err == nil && branch != "" {
			context = append(context, fmt.Sprintf("Branch: %s", branch))
		}

		recentCommits, err := gitOps.GetRecentCommits(3)
		if err == nil && len(recentCommits) > 0 {
			context = append(context, "Recent commits: "+strings.Join(recentCommits, ", "))
		}

		if verbose {
			fmt.Printf("📁 Current branch: %s\n", branch)
			fmt.Printf("📄 Staged changes: %d lines\n", strings.Count(diff, "\n"))
		}

		// Initialize AI client
		if verbose {
			fmt.Println("🧠 Generating commit message...")
		}

		aiClient, err := ai.NewClient(cfg)
		if err != nil {
			log.Fatalf("❌ Failed to initialize AI client: %v", err)
		}

		message, err := aiClient.GenerateCommitMessage(diff, context)
		if err != nil {
			log.Fatalf("❌ Error generating commit message: %v", err)
		}

		fmt.Println("\n✨ Generated commit message:")
		fmt.Printf("┌─%s─┐\n", strings.Repeat("─", len(message)))
		fmt.Printf("│ %s │\n", message)
		fmt.Printf("└─%s─┘\n", strings.Repeat("─", len(message)))

		if !skipConfirm {
			prompt := promptui.Prompt{
				Label:     "✅ Commit with this message?",
				IsConfirm: true,
			}

			_, err := prompt.Run()
			if err != nil {
				fmt.Println("Commit cancelled.")
				return
			}
		}

		err = gitOps.Commit(message)
		if err != nil {
			log.Fatalf("❌ Error committing: %v", err)
		}

		fmt.Println("🎉 Changes committed successfully!")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirmation and commit immediately")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output")
}
