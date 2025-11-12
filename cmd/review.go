/*
Copyright © 2025 Alexandro Cu alexandro.cuma@gmail.com
*/
package cmd

import (
	"fmt"
	"gommit/internal/config"
	"gommit/internal/git"
	"gommit/internal/helpers"
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
	Short: "Generate AI-powered PR descriptions from branch differences",
	Long: `Automatically create comprehensive PR descriptions using AI by analyzing changes between branches.

		This command compares your current branch with a base branch, analyzes the code differences,
		and generates a detailed PR description including summary, changes, and potential improvements.

		Features:
		• Compares changes between current branch and base branch
		• Generates comprehensive PR descriptions with AI
		• Includes code analysis and improvement suggestions
		• Works with your configured AI provider
		• Automatically detects git repository and branches

		Examples:
			gommit review                   # Compare with default base branch
			gommit review --base main       # Compare with main branch
			gommit review --base develop    # Compare with develop branch

		The generated PR description includes:
		• Overview of changes
		• Code analysis and impact
		• Potential issues or improvements
		• Ready-to-use PR description text`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("❌ Failed to load configuration: %v", err)
		}

		cfg.ValidateAIConfig()

		fmt.Println("🔍 Checking system requirements...")

		if utils.IsClipboardAvailable() {
			fmt.Println("✅ Clipboard support: Available")
		} else {
			fmt.Println("⚠️  Clipboard support: Not available")
			fmt.Printf("ℹ️  %s\n", utils.GetClipboardInfo())
		}

		// Check git
		gitOps := &git.RealGitOperations{}

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

		// Get diff between branches
		diff, err := gitOps.GetDiffBetweenBranches(baseBranch, currentBranch)
		if err != nil {
			log.Fatalf("❌ Failed to get diff: %v", err)
		}

		// Initialize AI client
		fmt.Println("🧠 Generating PR description with AI...")
		aiClient, err := ai.NewClient(cfg)
		if err != nil {
			log.Fatalf("❌ Failed to initialize AI client: %v", err)
		}

		// Generate PR description using template
		prReview, err := aiClient.GeneratePRReview(diff)
		if err != nil {
			log.Fatalf("❌ Error generating PR description: %v", err)
		}

		out, err := helpers.RenderMarkdown(prReview)
    if err != nil {
      out = prReview
    }

		// Display results
		fmt.Println("\n" + strings.Repeat("━", 60))
    fmt.Print(out)
		fmt.Println(strings.Repeat("━", 60))
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().StringVarP(&baseBranch, "base", "b", "", "Base branch to compare against (default: main/master/production)")
}
