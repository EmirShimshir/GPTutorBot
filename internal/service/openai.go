package service

import (
	"errors"
	"github.com/EmirShimshir/tasker-bot/internal/chatAI"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func (s *Service) ProcessTask(text string, chatID int64) (string, error) {
	var resErr error = nil
	ok, err := s.repo.Users.Exists(chatID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", NotAuthError
	}

	user, err := s.repo.Users.Get(chatID)
	if err != nil {
		return "", err
	}

	if !s.IsSubscriberUser(user) && user.Balance == 0 {
		return "", EmptyBalanceError
	}

	result, err := s.chatGpt.MakeRequest(text)
	if err != nil && !errors.Is(err, chatAI.ErrGptResult) && !errors.Is(err, chatAI.GptNewToken) {
		return "", err
	} else {
		resErr = err
	}

	if !s.IsSubscriberUser(user) {
		user.Balance--
	}

	err = s.repo.Users.Save(user)
	if err != nil {
		return "", err
	}

	log.WithFields(log.Fields{
		"chatID": chatID,
		"result": result,
	}).Info("ProcessTask")

	return result, resErr
}

func (s *Service) GetTokensDataAll() ([]byte, error) {
	log.Info("getTokensDataAll")
	file, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	defer os.Remove(file.Name())

	for _, string := range s.chatGpt.GetTokensAll() {
		_, err = file.WriteString(string)
		if err != nil {
			return nil, err
		}
	}

	fileData, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (s *Service) AddToken(token string) {
	s.chatGpt.AddToken(token)
}

func (s *Service) RemoveToken(token string) error {
	return s.chatGpt.RemoveToken(token)
}

func (s *Service) NextToken() {

	s.chatGpt.NextToken()
}
