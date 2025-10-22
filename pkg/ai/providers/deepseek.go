package providers

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type DeepSeekProvider struct {
	client *openai.Client
}

func NewDeepSeekProvider(apiKey string) *DeepSeekProvider {
	// DeepSeek uses OpenAI-compatible API
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.deepseek.com/v1" // DeepSeek API endpoint
	
	client := openai.NewClientWithConfig(config)
	return &DeepSeekProvider{
		client: client,
	}
}

func (p *DeepSeekProvider) Name() string {
	return "deepseek"
}

func (p *DeepSeekProvider) ValidateConfig(apiKey, model string) error {
	if apiKey == "" {
		return fmt.Errorf("DeepSeek API key is required")
	}
	return nil
}

func (p *DeepSeekProvider) GetDefaultModel() string {
	return "deepseek-chat"
}

func (p *DeepSeekProvider) CreateChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// Convert messages to OpenAI format
	var messages []openai.ChatCompletionMessage
	for _, msg := range req.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Use DeepSeek-specific model if needed
	model := req.Model
	if model == "gpt-4" || model == "gpt-3.5-turbo" {
		model = "deepseek-chat" // Map to DeepSeek equivalent
	}

	// Create completion request
	completionReq := openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: float32(req.Temperature),
		MaxTokens:   req.MaxTokens,
	}

	resp, err := p.client.CreateChatCompletion(ctx, completionReq)
	if err != nil {
		return nil, fmt.Errorf("DeepSeek API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no completion choices returned")
	}

	response := &ChatResponse{
		Content: resp.Choices[0].Message.Content,
	}
	response.Usage.PromptTokens = resp.Usage.PromptTokens
	response.Usage.CompletionTokens = resp.Usage.CompletionTokens
	response.Usage.TotalTokens = resp.Usage.TotalTokens

	return response, nil
}