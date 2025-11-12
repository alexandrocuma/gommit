package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// getConfigDir determines the appropriate config directory
func GetConfigDir() (string, error) {
	// If config was loaded, use its directory
	 configFile := viper.ConfigFileUsed();
	if configFile != "" {
		return filepath.Dir(configFile), nil
	}

	// Check local directory
	cwd, err := os.Getwd()
	if err == nil {
		localConfigPath := filepath.Join(cwd, "gommit", ".gommit.config.yaml")
		_, err := os.Stat(localConfigPath);
		if err == nil {
			return filepath.Dir(localConfigPath), nil
		}
	}

	// Default to home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %w", err)
	}
	return filepath.Join(home, ".gommit"), nil
}