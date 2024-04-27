package openAI

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	RoleContent string
	Token string
}

func NewOpenAI(roleContent, token string) *OpenAI {
	return &OpenAI{
		RoleContent: roleContent,
		Token: token,
	}
}

func (o *OpenAI) CreateChatCompletion(msg string) (string, error) {
	c := openai.NewClient(o.Token)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
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
