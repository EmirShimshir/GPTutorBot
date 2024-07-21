package openAI

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	BaseUrl     string
	RoleContent string
	Token       string
}

func NewOpenAI(baseUrl, roleContent, token string) *OpenAI {
	return &OpenAI{
		BaseUrl:     baseUrl,
		RoleContent: roleContent,
		Token:       token,
	}
}

func (o *OpenAI) CreateChatCompletion(msg string) (string, error) {
	config := openai.DefaultConfig(o.Token)
	config.BaseURL = o.BaseUrl
	c := openai.NewClientWithConfig(config)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: msg,
			},
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: o.RoleContent,
			},
		},
	}
	resp, err := c.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAI) NewToken(token string) {
	o.Token = token
}
