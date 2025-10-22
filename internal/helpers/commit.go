package helpers

import "gommit/internal/config"

func IndexToBool(i int) bool {
	return i == 0 // 0 = "Yes" = true
}

func GetCommitStyle(cfg *config.Config) string {
	style := "standard"
	if cfg.Commit.Conventional {
		style = "conventional commits"
	}
	if cfg.Commit.Emoji {
		style += " + emojis"
	}
	return style
}