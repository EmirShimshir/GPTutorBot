package deepSeek

import (
	"bytes"
	"encoding/json"
	"fmt"
"github.com/pkg/errors"
	"net/http"
)

var ErrStatus = errors.New("status")

type response struct {
	Choices           []choice `json:"choices"`
}

type choice struct {
	Message message `json:"message"`

}

type message struct {
	Role string `json:"role"`
	Content      string `json:"content"`
}

type request struct {
	Model             string   `json:"model"`
	Messages []message    `json:"messages"`

}

type DeepSeek struct {
	RoleContent string
	Token string
}

func NewDeepSeek(roleContent, token string) *DeepSeek {
	return &DeepSeek{
		RoleContent: roleContent,
		Token: token,
	}
}

func (d *DeepSeek) CreateChatCompletion(msg string) (string, error) {
	url := "https://api.deepseek.com/v1/chat/completions"
	data := request{
		Model: "deepseek-chat",
		Messages: []message {
			{
				Role: "user",
				Content: msg,
			},
			{
				Role: "system",
				Content: d.RoleContent,
			},
		},
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrap(ErrStatus, resp.Status)
	}

	var result response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return 	result.Choices[0].Message.Content, nil
}

func (d *DeepSeek) NewToken(token string) {
	d.Token = token
}