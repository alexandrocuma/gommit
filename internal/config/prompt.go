package config

type Prompt struct {
	Commit    string  `yaml:"commit" mapstructure:"commit"`
	Draft     string  `yaml:"draft" mapstructure:"draft"`
	Review    string  `yaml:"review" mapstructure:"review"`
}

func DefaultPromptConfig() *Prompt {
	cfg := &Prompt{}

	// Prompt defaults
	cfg.Commit = "commit.md"
	cfg.Draft  = "draft.md"
	cfg.Review = "review.md"

	return cfg
}
