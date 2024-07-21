package chatAI

import (
	"errors"
	"fmt"
	"github.com/EmirShimshir/tasker-bot/internal/chatAI/adapter/clientAI/openAI"
	"github.com/EmirShimshir/tasker-bot/internal/chatAI/adapter/queue/qList"
	"github.com/EmirShimshir/tasker-bot/internal/chatAI/port"
	"github.com/EmirShimshir/tasker-bot/internal/config"
	log "github.com/sirupsen/logrus"
)

var ErrGptResult = errors.New("error chatAI result")
var GptNewToken = errors.New("gpt chatAI token")

type ChatAI struct {
	client port.ClientAI
	queue  port.Queue
}

func NewChatAI(token string, cfg config.ChatAI) *ChatAI {
	t := port.NewToken(token, nil)
	q := qList.NewQList()
	q.Add(t)
	return &ChatAI{
		client: openAI.NewOpenAI(cfg.BaseUrl, cfg.RoleContent, q.Get().Value),
		queue:  q,
	}
}

func (c *ChatAI) GetTokensAll() []string {
	res := make([]string, 0, 1)
	tokens := c.queue.GetAll()
	for i, el := range tokens {
		errRes := "nil"
		if el.Err != nil {
			errRes = el.Err.Error()
		}
		res = append(res, fmt.Sprintf("%d) TOKEN: %s; ERROR: %s; CNTREQ: %d\n", len(tokens)-i, el.Value, errRes, el.CntReq))
	}

	return res
}

func (c *ChatAI) AddToken(token string) {

	t := port.NewToken(token, nil)
	c.queue.Add(t)
	c.client.NewToken(c.queue.Get().Value)
}

func (c *ChatAI) RemoveToken(token string) error {
	err := c.queue.Remove(token)
	if err != nil {
		return err
	}

	c.client.NewToken(c.queue.Get().Value)
	return nil
}

func (c *ChatAI) NextToken() {

	c.queue.Next()
	c.client.NewToken(c.queue.Get().Value)
}

func (c *ChatAI) MakeRequest(message string) (string, error) {
	res, err := c.client.CreateChatCompletion(message)
	c.queue.Get().CntReq++
	if err != nil {
		token := c.queue.Get()
		token.Err = err

		log.WithFields(log.Fields{
			"token": token.Value,
			"err":   token.Err,
		}).Info("ChatGPT token error INFO")

		c.queue.Next()
		for token.Value != c.queue.Get().Value {
			c.client.NewToken(c.queue.Get().Value)
			res, err = c.client.CreateChatCompletion(message)
			c.queue.Get().CntReq++
			if err != nil {
				token := c.queue.Get()
				token.Err = err

				log.WithFields(log.Fields{
					"token": token.Value,
					"err":   token.Err,
				}).Info("ChatGPT token error INFO")
			} else {
				return res, GptNewToken
			}
			c.queue.Next()
		}
		log.WithFields(log.Fields{
			"token": token.Value,
			"err":   token.Err,
		}).Info("ChatGPT token error ALERT")

		return "", ErrGptResult

	}

	return res, nil

}
