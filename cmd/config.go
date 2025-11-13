/*
Copyright © 2025 Alexandro Cu alexandro.cuma@gmail.com
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/alexandrocuma/gommit/internal/config"
	"github.com/alexandrocuma/gommit/internal/helpers"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display current configuration settings",
	Long: `Shows your active Gommit configuration including AI provider settings and config file location.

			Features:
			• Displays AI provider, model, and parameters
			• Shows masked API key for security
			• Reveals configuration file path
			• Helps verify and debug settings
			• Validates configuration loading

			Examples:
				gommit config
				gommit config show

			The output includes:
			• AI provider configuration
			• Model settings (temperature, max tokens)
			• Masked API key (showing last 4 characters)
			• Config file location on disk`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("❌ Failed to load configuration: %v", err)
		}

		fmt.Println("📋 Current Gommit Configuration:")
		fmt.Println("────────────────────────────────")

		fmt.Printf("\n🤖 AI Settings:\n")
		fmt.Printf("  Provider:    %s\n", cfg.AI.Provider)
		fmt.Printf("  Model:       %s\n", cfg.AI.Model)
		fmt.Printf("  Temperature: %.1f\n", cfg.AI.Temperature)
		fmt.Printf("  Max Tokens:  %d\n", cfg.AI.MaxTokens)
		fmt.Printf("  API Key:     %s\n", helpers.MaskAPIKey(cfg.AI.APIKey))

		fmt.Printf("\n📁 Config file: %s\n", helpers.GetConfigPath())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
