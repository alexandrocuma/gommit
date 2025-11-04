package providers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

type AnthropicProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicRequest struct {
	Model       string             `json:"model"`
	Messages    []AnthropicMessage `json:"messages"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature float64            `json:"temperature,omitempty"`
}

type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func NewAnthropicProvider(apiKey string) *AnthropicProvider {
	return &AnthropicProvider{
		apiKey:     apiKey,
		baseURL:    "https://api.anthropic.com/v1/messages",
		httpClient: &http.Client{},
	}
}

func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

func (p *AnthropicProvider) ValidateConfig(apiKey, model string) error {
	if apiKey == "" {
		return fmt.Errorf("anthropic API key is required")
	}
	if !strings.HasPrefix(apiKey, "sk-ant-") {
		return fmt.Errorf("invalid Anthropic API key format")
	}
	return nil
}

func (p *AnthropicProvider) GetDefaultModel() string {
	return "claude-3-sonnet-20240229"
}

func (p *AnthropicProvider) CreateChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
    client := anthropic.NewClient(
			option.WithAPIKey(p.apiKey),
		)
    
    // Separate system messages from conversation messages
    var systemContent string
    var conversationMessages []anthropic.MessageParam
    
    for _, msg := range req.Messages {
        switch msg.Role {
        case "system":
            // Combine all system messages
            if systemContent != "" {
                systemContent += "\n"
            }
            systemContent += msg.Content
        case "user":
            conversationMessages = append(conversationMessages, 
                anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content)))
        case "assistant":
            conversationMessages = append(conversationMessages, 
                anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)))
        }
    }
    
    // Build the request parameters
    params := anthropic.MessageNewParams{
        Model:     anthropic.Model(req.Model),
        MaxTokens: int64(req.MaxTokens),
        Messages:  conversationMessages,
				Temperature: param.Opt[float64]{Value: req.Temperature},
    }
    
    // Add system message if present
    if systemContent != "" {
        params.System = []anthropic.TextBlockParam{
            {Type: "text", Text: systemContent},
        }
    }
    
    // Make the API call
    response, err := client.Messages.New(ctx, params)
    if err != nil {
        return nil, fmt.Errorf("anthropic API error: %w", err)
    }
    
    return &ChatResponse{
        Content: response.Content[0].Text,
        Usage: struct {
            PromptTokens     int `json:"prompt_tokens"`
            CompletionTokens int `json:"completion_tokens"`
            TotalTokens      int `json:"total_tokens"`
        }{
            PromptTokens:     int(response.Usage.InputTokens),
            CompletionTokens: int(response.Usage.OutputTokens),
            TotalTokens:      int(response.Usage.InputTokens + response.Usage.OutputTokens),
        },
    }, nil
}