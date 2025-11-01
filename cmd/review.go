/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gommit/internal/config"
	"gommit/internal/git"
	"gommit/pkg/ai"
	"gommit/pkg/utils"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("❌ Failed to load configuration: %v", err)
		}

		// Check if AI is configured
		if cfg.AI.APIKey == "" {
			fmt.Println("❌ AI API key not configured.")
			fmt.Println("Please run 'gitai init' to set up your configuration.")
			os.Exit(1)
		}
		fmt.Println("🔍 Checking system requirements...")

		if utils.IsClipboardAvailable() {
			fmt.Println("✅ Clipboard support: Available")
		} else {
			fmt.Println("⚠️  Clipboard support: Not available")
			fmt.Printf("ℹ️  %s\n", utils.GetClipboardInfo())
		}

		// Check git
		gitOps := &git.RealGitOperations{}
		if gitOps.IsGitRepository() {
			fmt.Println("✅ Git repository: Detected")
		} else {
			fmt.Println("⚠️  Git repository: Not detected (will need to be in a git repo to use commit/PR features)")
		}
		// Get current branch
		currentBranch, err := gitOps.GetCurrentBranch()
		if err != nil {
			log.Fatalf("❌ Failed to get current branch: %v", err)
		}

		// Set default base branch if not provided
		if baseBranch == "" {
			baseBranch = getDefaultBaseBranch(gitOps)
		}

		fmt.Printf("📊 Comparing changes from '%s' to '%s'...\n", currentBranch, baseBranch)

		// Get diff between branches
		diff, err := gitOps.GetDiffBetweenBranches(baseBranch, currentBranch)
		if err != nil {
			log.Fatalf("❌ Failed to get diff: %v", err)
		}

		if diff == "" {
			fmt.Println("❌ No changes found between branches.")
			fmt.Println("   Make sure you have committed your changes and the branches are different.")
			os.Exit(1)
		}

		// Initialize AI client
		fmt.Println("🧠 Generating PR description with AI...")
		aiClient, err := ai.NewClient(cfg)
		if err != nil {
			log.Fatalf("❌ Failed to initialize AI client: %v", err)
		}

		// Generate PR description using template
		prReview, err := aiClient.GenerateReview(diff)
		if err != nil {
			log.Fatalf("❌ Error generating PR description: %v", err)
		}

		// Display results
		fmt.Println("\n" + strings.Repeat("━", 60))
		fmt.Println("📋 PR Review generated")
		fmt.Println(strings.Repeat("━", 60))
		fmt.Println(prReview)
		fmt.Println(strings.Repeat("━", 60))
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().StringVarP(&baseBranch, "base", "b", "", "Base branch to compare against (default: main/master)")
	reviewCmd.Flags().StringVarP(&templateFile, "template", "t", "default", "Template name or path to template file")
	reviewCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file to save PR description")
	reviewCmd.Flags().BoolVar(&skipReview, "skip-review", false, "Skip interactive review and editing")
	reviewCmd.Flags().BoolVarP(&copyToClipboard, "clipboard", "c", false, "Copy PR description to clipboard")
}