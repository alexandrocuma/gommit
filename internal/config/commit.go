package config

type Commit struct {
	Conventional bool   `yaml:"conventional" mapstructure:"conventional"`
	Emoji        bool   `yaml:"emoji" mapstructure:"emoji"`
	Language     string `yaml:"language" mapstructure:"language"`
}

func DefaultCommitConfig() *Commit {
	cfg := &Commit{}

	// AI defaults
	cfg.Conventional = true
	cfg.Emoji = false
	cfg.Language = "english"

	return cfg
}
