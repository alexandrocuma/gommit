/*
Copyright © 2025 Alexandro Cu alexandro.cuma@gmail.com
*/
package cmd

import (
	"fmt"
	"gommit/internal/config"
	"gommit/internal/helpers"
	"log"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		fmt.Printf("\n📁 Config files: %s\n", helpers.GetConfigPath())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
