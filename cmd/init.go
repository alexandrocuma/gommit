/*
Copyright © 2025 Alexandro Cu alexandro.cuma@gmail.com
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/alexandrocuma/gommit/internal/config"
	"github.com/alexandrocuma/gommit/internal/helpers"
	"github.com/alexandrocuma/gommit/pkg/interactive"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration interactively",
	Long: `Runs an interactive setup wizard to create or overwrite your Gommit configuration file.

			Features:
			• Interactive prompt-driven configuration
			• Guides through AI provider selection
			• Securely stores API credentials
			• Warns before overwriting existing config
			• Automatically saves to correct location
			• Validates configuration before saving

			Examples:
				gommit init

			Setup process includes:
			• AI provider and model selection
			• API key configuration
			• Parameter tuning (temperature, tokens)
			• Config file creation in user directory
			• Next steps guidance for using gommit`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.ConfigExists() {
			fmt.Println("⚠️  Configuration file already exists!")
			prompt := promptui.Prompt{
				Label:     "Do you want to overwrite it",
				IsConfirm: true,
			}

			_, err := prompt.Run()
			if err != nil {
				fmt.Println("Init cancelled.")
				return
			}
		}

		cfg, err := interactive.RunSetup()
		if err != nil {
			log.Fatalf("❌ Init failed: %v", err)
		}

		err = config.SaveConfig(cfg)
		if err != nil {
			log.Fatalf("❌ Failed to save configuration: %v", err)
		}

		configPath := helpers.GetConfigPath()
		fmt.Printf("✅ Configuration saved to: %s\n", configPath)
		fmt.Println("\n🎉 Setup complete! You can now use gommit CLI.")
		fmt.Println("\nNext steps:")
		fmt.Println("  gommit     		 # Generate commit messages")
		fmt.Println("  gommit draft    # Generate a PR description")
		fmt.Println("  gommit review   # Generate a PR review")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
