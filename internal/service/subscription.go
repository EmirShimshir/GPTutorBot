package service

import (
	"github.com/EmirShimshir/tasker-bot/internal/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

func (s *Service) UpdateSubscription(chatID int64, countBought int64) error {
	log.WithFields(log.Fields{
		"chatID":      chatID,
		"countBought": countBought,
	}).Info("UpdateSubscription")

	ok, err := s.repo.Users.Exists(chatID)
	if err != nil {
		return err
	}
	if !ok {
		return NotAuthError
	}

	user, err := s.repo.Users.Get(chatID)
	if err != nil {
		return err
	}

	now := time.Now()

	diff := now.Sub(user.DateSub)

	if diff > 0 {
		user.DateSub = now.AddDate(0, int(countBought), 0)
	} else {
		user.DateSub = user.DateSub.AddDate(0, int(countBought), 0)
	}

	return s.repo.Users.Save(user)
}

func (s *Service) GetSubscribeDateEnd(chatID int64) (string, error) {
	log.WithFields(log.Fields{
		"chatID": chatID,
	}).Info("GetSubscribeDateEnd")

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

	return user.DateSub.Format("02.01.2006"), nil
}

func (s *Service) IsSubscriber(chatID int64) (bool, error) {
	ok, err := s.repo.Users.Exists(chatID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, NotAuthError
	}

	user, err := s.repo.Users.Get(chatID)
	if err != nil {
		return false, err
	}

	now := time.Now()

	diff := now.Sub(user.DateSub)

	if diff > 0 {
		return false, nil
	}

	return true, nil
}

func (s *Service) IsSubscriberUser(user *domain.User) bool {
	now := time.Now()
	diff := now.Sub(user.DateSub)

	if diff > 0 {
		return false
	}

	return true
}
