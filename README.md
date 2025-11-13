## Gommit

AI-powered Git companion

![Go Version](https://img.shields.io/badge/go-%3E%3D%201.21-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

Gommit is a command-line tool that leverages AI to automate your Git workflow. It generates intelligent commit messages, creates comprehensive PR descriptions, and helps you write better documentation for your code changes.

## Features

- 🤖 AI-Agnostic: Integrates with OpenAI, Anthropic, and other providers
- 💬 Smart Commit Messages: Analyzes your staged changes and suggests meaningful commits
- 📋 PR Description Generation: Automatically creates detailed PR descriptions from branch differences
- 📝 PR Review Generation: Automatically creates detailed PR review from branch differences
- 🎯 Template Support: Customize output with Go templates
- 🔒 Secure: API keys are stored locally and masked in output
- ⚡ Fast: Works directly with Git operations for quick analysis

## Installation

### Using Go Install

```bash
go install github.com/alexandrocuma/gommit@latest
```

### Build from Source

```bash
git clone https://github.com/alexandrocuma/gommit.git
cd gommit
go build -o gommit main.go
sudo mv gommit /usr/local/bin/
```

## Quick Start

### Initialize Configuration

Run the interactive setup wizard to configure your AI provider:

```bash
gommit init
```

This will guide you through:

- Selecting an AI provider (OpenAI, Anthropic, etc.)
- Entering your API key
- Choosing a model
- Configuring generation parameters

---

### View Your Config

Verify your configuration at any time:

```bash
gommit config
```

---

### Generate a PR Description

Create a comprehensive PR description from your branch changes:

```bash
gommit draft --base main
```

## 🧩 Command Reference

```bash
gommit init
```

Initialize your Gommit configuration interactively.

**Features:**

- Interactive prompt-driven configuration
- Guides through AI provider selection
- Securely stores API credentials
- Warns before overwriting existing config
- Automatically saves to correct location

**Examples:**

- gommit init — Run first-time setup

---

```bash
gommit draft
```

Generate PR descriptions from branch differences.

**Features:**

- Compares current branch with any base branch
- Analyzes commits, diff stats, and code changes
- Uses customizable templates for structure
- Generates intelligent PR titles from branch names
- Supports saving to file or clipboard
- Works with your configured AI provider

**Examples:**

- gommit draft — Compare with default base branch
- gommit draft --base main — Compare with main branch
- gommit draft --base develop — Compare with develop branch
- gommit draft --title "My changes" — Use custom PR title
- gommit draft --output pr.md — Save to file

---

```bash
gommit config
```

Display current configuration settings.

**Features:**

- Displays AI provider, model, and parameters
- Shows masked API key for security
- Reveals configuration file path
- Helps verify and debug settings
- Validates configuration loading

**Examples:**

- gommit config

**The output includes:**

- AI provider configuration
- Model settings (temperature, max tokens)
- Masked API key (showing last 4 characters)
- Config file location on disk

---

## ⚙️ Configuration

Configuration is stored in:

- **Linux/macOS:** `~/.config/gommit/config.yaml`
- **Windows:** `%APPDATA%\gommit\config.yaml`

**Example configuration:**

```yaml
ai:
  provider: openai
  model: gpt-4o-mini
  temperature: 0.7
  max_tokens: 2048
  api_key: sk-...your-key-here
```

## 📁 Project Structure

```bash
.
├── cmd/           # CLI commands
├── internal/      # AI provider, Configurations, Git integrations
├── pkg/           # shared packages for integrations
├── main.go        # Entry point
└── README.md      # This file
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
