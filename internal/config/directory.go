package config

type Directory struct {
	Prompts    string `yaml:"prompts" mapstructure:"prompts"`
	Templates  string `yaml:"templates" mapstructure:"templates"`
}

func DefaultDirectoryConfig() *Directory {
	cfg := &Directory{}

	// Directory defaults
	cfg.Prompts   = "~/.gommit"
	cfg.Templates = "./templates" // Relative to current working directory

	return cfg
}
