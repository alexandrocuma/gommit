package helpers

import "os"

// Helper functions
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.gommit.config.yaml"
	}
	return home + "/.gommit.config.yaml"
}

func MaskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "***"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}
