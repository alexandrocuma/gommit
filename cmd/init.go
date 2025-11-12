/*
Copyright © 2025 Alexandro Cu alexandro.cuma@gmail.com
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
		fmt.Println("  gommit     		 # Generate AI-powered commit messages")
		fmt.Println("  gommit draft    # Generate a PR description")
		fmt.Println("  gommit review   # Generate a PR review")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
