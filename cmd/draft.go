/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"gommit/internal/config"
	"gommit/internal/git"
	"gommit/pkg/ai"
	"gommit/pkg/template"
	"gommit/pkg/utils"
	"log"
	"os"
	"os/exec"
	"strings"

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
			prTitle = generatePRTitle(currentBranch, commits)
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

		// Also update the clipboard usage in the main PR function:
		if copyToClipboard {
			fullPRContent := fmt.Sprintf("# %s\n\n%s", prTitle, prDescription)
			if err := copyToClipboardUtil(fullPRContent); err != nil {
				fmt.Printf("⚠️  Failed to copy to clipboard: %v\n", err)
				fmt.Printf("ℹ️  %s\n", utils.GetClipboardInfo())
			} else {
				fmt.Println("📋 PR description copied to clipboard!")
			}
		}

		// Interactive review
		if !skipReview {
			fmt.Print("\n🤔 Would you like to edit the PR description? [y/N]: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			response := strings.ToLower(strings.TrimSpace(scanner.Text()))

			if response == "y" || response == "yes" {
				editedContent, err := openInEditor(fmt.Sprintf("# %s\n\n%s", prTitle, prDescription))
				if err != nil {
					fmt.Printf("⚠️  Failed to open editor: %v\n", err)
				} else {
					fmt.Println("✅ PR description updated with your edits!")
					if outputFile != "" {
						if err := os.WriteFile(outputFile, []byte(editedContent), 0644); err != nil {
							fmt.Printf("⚠️  Failed to update output file: %v\n", err)
						}
					}
				}
			}
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

// Helper functions
func getDefaultBaseBranch(gitOps git.GitOperations) string {
	// Try common base branch names
	possibleBranches := []string{"main", "master", "develop"}

	for _, branch := range possibleBranches {
		if gitOps.BranchExists(branch) {
			return branch
		}
	}

	// Fallback to main
	return "main"
}

func generatePRTitle(currentBranch string, commits []string) string {
	// Clean up branch name for PR title
	title := strings.TrimPrefix(currentBranch, "feature/")
	title = strings.TrimPrefix(title, "feat/")
	title = strings.TrimPrefix(title, "fix/")
	title = strings.TrimPrefix(title, "bugfix/")
	title = strings.TrimPrefix(title, "hotfix/")

	// Convert kebab-case or snake_case to Title Case
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")
	title = strings.Title(title)

	return title
}

// Replace the copyToClipboardUtil function with:
func copyToClipboardUtil(content string) error {
	if !utils.IsClipboardAvailable() {
		return fmt.Errorf("clipboard not available on this system")
	}

	if err := utils.CopyToClipboard(content); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	return nil
}

func openInEditor(content string) (string, error) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "pr-description-*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	// Write content to temp file
	if _, err := tmpFile.WriteString(content); err != nil {
		return "", err
	}
	tmpFile.Close()

	// Get editor from environment or use default
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // default to vim
	}

	// Open editor
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Read edited content
	editedContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}

	return string(editedContent), nil
}
