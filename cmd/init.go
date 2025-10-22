/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gommit/internal/config"
	"gommit/internal/helpers"
	"gommit/pkg/interactive"
	"log"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if config already exists
		if config.ConfigExists() {
			fmt.Println("⚠️  Configuration file already exists!")
			fmt.Print("Do you want to overwrite it? [y/N]: ")

			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "Y" {
				fmt.Println("Init cancelled.")
				return
			}
		}

		// Run interactive setup
		cfg, err := interactive.RunSetup()
		if err != nil {
			log.Fatalf("❌ Init failed: %v", err)
		}

		// Save configuration
		err = config.SaveConfig(cfg)
		if err != nil {
			log.Fatalf("❌ Failed to save configuration: %v", err)
		}

		configPath := helpers.GetConfigPath()
		fmt.Printf("✅ Configuration saved to: %s\n", configPath)
		fmt.Println("\n🎉 Setup complete! You can now use GitAI CLI.")
		// fmt.Println("\nNext steps:")
		// fmt.Println("  gitai commit    # Generate AI-powered commit messages")
		// fmt.Println("  gitai pr        # Create PR descriptions")
		// fmt.Println("  gitai config    # Manage your configuration")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
