package providers

import "context"

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Content string `json:"content"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Provider defines the interface for AI providers
type Provider interface {
	// Name returns the provider name
	Name() string

	// CreateChatCompletion sends a chat completion request
	CreateChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// ValidateConfig validates the provider configuration
	ValidateConfig(apiKey, model string) error

	// GetDefaultModel returns the default model for this provider
	GetDefaultModel() string
}
