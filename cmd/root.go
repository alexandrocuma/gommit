/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"gommit/internal/config"
	"gommit/internal/git"
	"gommit/internal/helpers"
	"gommit/pkg/ai"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	skipConfirm bool
	verbose     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gommit",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
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

		if verbose {
			fmt.Printf("🤖 Using AI provider: %s\n", cfg.AI.Provider)
			fmt.Printf("📝 Commit style: %s\n", helpers.GetCommitStyle(cfg))
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
			baseBranch = getDefaultBaseBranch(gitOps)
		}

		fmt.Printf("📊 Comparing changes from '%s' to '%s'...\n", currentBranch, baseBranch)
		// Get staged diff
		if verbose {
			fmt.Println("📊 Analyzing staged changes...")
		}

		diff, err := gitOps.GetStagedDiff()
		if err != nil {
			log.Fatalf("❌ Error getting git diff: %v", err)
		}

		if diff == "" {
			fmt.Println("❌ No staged changes found.")
			fmt.Println("   Please stage your changes first: git add <files>")
			os.Exit(1)
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
			fmt.Println("🧠 Generating commit message with AI...")
		}

		aiClient, err := ai.NewClient(cfg)
		if err != nil {
			log.Fatalf("❌ Failed to initialize AI client: %v", err)
		}

		message, err := aiClient.GenerateCommitMessage(diff, context, &cfg.Commit)
		if err != nil {
			log.Fatalf("❌ Error generating commit message: %v", err)
		}
		// Display the generated message
		fmt.Println("\n✨ Generated commit message:")
		fmt.Println("┌─────────────────────────────────────────────────────────┐")
		fmt.Printf("│ %s\n", message)
		fmt.Println("└─────────────────────────────────────────────────────────┘")

		// Ask for confirmation unless skipped
		if !skipConfirm {
			fmt.Print("\n✅ Commit with this message? [Y/n]: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			response := strings.ToLower(strings.TrimSpace(scanner.Text()))

			if response == "n" || response == "no" {
				fmt.Println("Commit cancelled.")
				return
			}
		}

		// Perform the commit
		if err := gitOps.Commit(message); err != nil {
			log.Fatalf("❌ Error committing: %v", err)
		}

		fmt.Println("🎉 Changes committed successfully!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.convmit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirmation and commit immediately")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output")
}
