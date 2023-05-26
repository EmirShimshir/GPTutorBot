package openai

import (
	"context"
	"fmt"
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"github.com/sashabaranov/go-openai"
)

type ChatGpt struct {
	client       *openai.Client
	contentStart string
}

func NewChatGpt(token string, cfg config.ChatGpt) *ChatGpt {
	return &ChatGpt{
		client:       openai.NewClient(token),
		contentStart: cfg.ContentStart,
	}
}

func (c *ChatGpt) MakeRequest(message string) (string, error) {
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("%s:\n%s", c.contentStart, message),
			},
		},
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil

}
