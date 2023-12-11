package service

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func (s *Service) IsGift(text string) bool {
	return text == s.promo.Gift
}

func (s *Service) ActivateGift(chatID int64) error {
	log.WithFields(log.Fields{
		"chatID": chatID,
	}).Info("ActivateGift")

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

	if user.UsedPromo == true {
		return PromoUsedError
	}
	user.UsedPromo = true

	now := time.Now()

	diff := now.Sub(user.DateSub)

	countGift := 7

	if diff > 0 {
		user.DateSub = now.AddDate(0, 0, countGift)
	} else {
		user.DateSub = user.DateSub.AddDate(0, 0, countGift)
	}

	return s.repo.Users.Save(user)
}
