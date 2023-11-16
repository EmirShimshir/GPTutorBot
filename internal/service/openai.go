package service

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (s *Service) ProcessTask(text string, chatID int64) (string, error) {
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
	if err != nil {
		return "", errors.Wrap(err, "MakeRequest\n"+result)
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

	return result, nil
}
