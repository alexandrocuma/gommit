/*
Copyright © 2025 Alexandro Cu alexandro.cuma@gmail.com
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/alexandrocuma/gommit/internal/config"
	"github.com/alexandrocuma/gommit/pkg/directory"
	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
)

// reviewCmd represents the review command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Generate PR descriptions from branch differences",
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

		files, err := directory.WalkFiles("./", func(path string, info fs.FileInfo) bool {
			return strings.HasSuffix(strings.ToLower(path), ".go")
		})
		if err != nil {
			fmt.Println("Failed to walk files")
			return
		}
		
		prompt := promptui.Select{
			Label: "Do you want to overwrite it",
			Items: files,
		}

		_, value, err := prompt.Run()
		if err != nil {
			fmt.Println("Init cancelled.")
			return
		}

		content, err := directory.LoadFileString(value)
		if err != nil {
			fmt.Println("Failed to load file content.")
			return
		}
		fmt.Println(content)
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
