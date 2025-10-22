package providers

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	client := openai.NewClient(apiKey)
	return &OpenAIProvider{
		client: client,
	}
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) ValidateConfig(apiKey, model string) error {
	if apiKey == "" {
		return fmt.Errorf("OpenAI API key is required")
	}
	if !strings.HasPrefix(apiKey, "sk-") {
		return fmt.Errorf("invalid OpenAI API key format")
	}
	return nil
}

func (p *OpenAIProvider) GetDefaultModel() string {
	return "gpt-4"
}

func (p *OpenAIProvider) CreateChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// Convert messages to OpenAI format
	var messages []openai.ChatCompletionMessage
	for _, msg := range req.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Create completion request
	completionReq := openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		Temperature: float32(req.Temperature),
		MaxTokens:   req.MaxTokens,
	}

	resp, err := p.client.CreateChatCompletion(ctx, completionReq)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
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