package openai

import (
	"context"
	"errors"
	"fmt"
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"github.com/EmirShimshir/tasker-bot/internal/queue"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
)

var ErrTokenFound = errors.New("value not found in queue")
var ErrGptResult = errors.New("error gpt result")
var GptNewToken = errors.New("gpt new token")

type Token struct {
	Value string
	Err   error
	CntReq int64
}

func NewToken(value string, err error) *Token {
	return &Token{
		Value: value,
		Err:   err,
		CntReq: 0,
	}
}

type ChatGpt struct {
	client       *openai.Client
	contentStart string
	queue *queue.Queue
}

func NewChatGpt(token string, cfg config.ChatGpt) *ChatGpt {
	t := NewToken(token, nil)
	q := queue.NewQueue()
	q.Add(t)
	return &ChatGpt{
		client:       openai.NewClient(q.Get().(*Token).Value),
		contentStart: cfg.ContentStart,
		queue: q,
	}
}

func (c *ChatGpt) GetTokensAll() []string {
	res := make([]string, 0, 1)

	tokens := c.queue.GetAll()
	for i, el := range tokens {
		errRes := "nil"
		if el.(*Token).Err != nil {
			errRes = el.(*Token).Err.Error()
		}
		res = append(res, fmt.Sprintf("%d) TOKEN: %s; ERROR: %s; CNTREQ: %d\n", len(tokens) - i, el.(*Token).Value, errRes, el.(*Token).CntReq))
	}

	return res
}

func (c *ChatGpt) AddToken(token string) {

	t := NewToken(token, nil)
	c.queue.Add(t)
	c.client = openai.NewClient(c.queue.Get().(*Token).Value)
}

func (c *ChatGpt) RemoveToken(token string) error {
	for e := c.queue.List.Front(); e != nil; e = e.Next() {
		if e.Value.(*Token).Value == token {
			c.queue.List.Remove(e)
			c.client = openai.NewClient(c.queue.Get().(*Token).Value)
			return nil
		}
	}
	return ErrTokenFound
}

func (c *ChatGpt) NextToken() {

	c.queue.Next()
	c.client = openai.NewClient(c.queue.Get().(*Token).Value)
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
	c.queue.Get().(*Token).CntReq++
	if err != nil {
		token := c.queue.Get().(*Token)
		token.Err = err

		log.WithFields(log.Fields{
			"token": token.Value,
			"err": token.Err,
		}).Info("ChatGPT token error INFO")

		c.queue.Next()
		for token.Value != c.queue.Get().(*Token).Value {
			c.client = openai.NewClient(c.queue.Get().(*Token).Value)
			resp, err = c.client.CreateChatCompletion(ctx, req)
			c.queue.Get().(*Token).CntReq++
			if err != nil {
				token := c.queue.Get().(*Token)
				token.Err = err

				log.WithFields(log.Fields{
					"token": token.Value,
					"err":   token.Err,
				}).Info("ChatGPT token error INFO")
			} else {
				return resp.Choices[0].Message.Content, GptNewToken
			}
			c.queue.Next()
		}
		log.WithFields(log.Fields{
			"token": token.Value,
			"err": token.Err,
		}).Info("ChatGPT token error ALERT")

		return "", ErrGptResult

	}

	return resp.Choices[0].Message.Content, nil

}
