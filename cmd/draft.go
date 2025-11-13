/*
Copyright © 2025 Alexandro Cu alexandro.cuma@gmail.com
*/
package cmd

import (
	"fmt"
	"gommit/internal/config"
	"gommit/internal/git"
	"gommit/pkg/ai"
	"gommit/pkg/template"
	"gommit/pkg/utils"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	baseBranch      string
	templateFile    string
	outputFile      string
	prTitle         string
	skipReview      bool
	copyToClipboard bool
)

// draftCmd represents the draft command
var draftCmd = &cobra.Command{
	Use:   "draft",
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

		cfg.ValidateAIConfig()

		fmt.Println("🔍 Checking system requirements...")

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

		// Get commit history
		commits, err := gitOps.GetCommitsBetweenBranches(baseBranch, currentBranch)
		if err != nil {
			log.Fatalf("❌ Failed to get commit history: %v", err)
		}

		// Get diff stats
		diffStats, err := gitOps.GetDiffStatsBetweenBranches(baseBranch, currentBranch)
		if err != nil {
			log.Fatalf("❌ Failed to get diff stats: %v", err)
		}

		// Load template
		templateContent, err := template.LoadTemplate(templateFile)
		if err != nil {
			log.Fatalf("❌ Failed to load template: %v", err)
		}

		fmt.Printf("📝 Using template: %s\n", templateFile)
		fmt.Printf("📄 Found %d commits with %d lines changed\n", len(commits), strings.Count(diff, "\n"))

		// Generate PR title if not provided
		if prTitle == "" {
			prTitle = generatePRTitle(currentBranch)
		}

		// Initialize AI client
		fmt.Println("🧠 Generating PR description with AI...")
		aiClient, err := ai.NewClient(cfg)
		if err != nil {
			log.Fatalf("❌ Failed to initialize AI client: %v", err)
		}

		// Generate PR description using template
		prDescription, err := aiClient.GeneratePRDescriptionWithTemplate(prTitle, commits, diff, diffStats, templateContent)
		if err != nil {
			log.Fatalf("❌ Error generating PR description: %v", err)
		}

		// Display results
		fmt.Println("\n" + strings.Repeat("━", 60))
		fmt.Println("📋 PR DESCRIPTION GENERATED")
		fmt.Println(strings.Repeat("━", 60))
		fmt.Printf("📌 Title: %s\n\n", prTitle)
		fmt.Println(prDescription)
		fmt.Println(strings.Repeat("━", 60))

		// Handle output options
		if outputFile != "" {
			fullPRContent := fmt.Sprintf("# %s\n\n%s", prTitle, prDescription)
			if err := os.WriteFile(outputFile, []byte(fullPRContent), 0644); err != nil {
				log.Fatalf("❌ Failed to write output file: %v", err)
			}
			fmt.Printf("💾 PR description saved to: %s\n", outputFile)
		}

		prompt := promptui.Prompt{
			Label:     "📄 Copy to clipboard",
			IsConfirm: true,
		}
		_, err = prompt.Run()
		
		if err != nil {
			fmt.Println("\n🎉 PR description ready!")
			return 
		}
		
		fullPRContent := fmt.Sprintf("# %s\n\n%s", prTitle, prDescription)
		err = utils.CopyToClipboardUtil(fullPRContent)
		if err == nil {
			fmt.Println("📋 PR description copied to clipboard!")
		}

		fmt.Println("\n🎉 PR description ready!")
	},
}

func init() {
	rootCmd.AddCommand(draftCmd)

	draftCmd.Flags().StringVarP(&baseBranch, "base", "b", "", "Base branch to compare against (default: main/master)")
	draftCmd.Flags().StringVarP(&templateFile, "template", "t", "default", "Template name or path to template file")
	draftCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file to save PR description")
	draftCmd.Flags().StringVarP(&prTitle, "title", "T", "", "PR title (default: auto-generated from branch name)")
	draftCmd.Flags().BoolVar(&skipReview, "skip-review", false, "Skip interactive review and editing")
	draftCmd.Flags().BoolVarP(&copyToClipboard, "clipboard", "c", false, "Copy PR description to clipboard")
}

func generatePRTitle(currentBranch string) string {
	// Clean up branch name for PR title
	title := strings.TrimPrefix(currentBranch, "feature/")
	title = strings.TrimPrefix(title, "feat/")
	title = strings.TrimPrefix(title, "fix/")
	title = strings.TrimPrefix(title, "bugfix/")
	title = strings.TrimPrefix(title, "hotfix/")

	// Convert kebab-case or snake_case to Title Case
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")

	return title
}

