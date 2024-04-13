package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (s *Service) CreateRef(chatID int64) (string, error) {
	log.WithFields(log.Fields{
		"chatID": chatID,
	}).Info("CreateRef")

	promo := fmt.Sprintf("https://t.me/%s?start=%s%d", s.promo.BotName, s.promo.Start, chatID)

	return promo, nil
}

func (s *Service) IsRef(text string) bool {

	return strings.HasPrefix(text, s.promo.Start)
}

func (s *Service) ActivateRef(text string, chatID int64) (int64, string, error) {
	log.WithFields(log.Fields{
		"chatID": chatID,
	}).Info("ActivatePromo")

	chatIdAuthor, err := strconv.ParseInt(strings.TrimPrefix(text, s.promo.Start), 10, 64)
	if err != nil {
		return 0, "", PromoError
	}

	ok, err := s.repo.Users.Exists(chatID)
	if err != nil {
		return 0, "", err
	}
	if !ok {
		return 0, "", NotAuthError
	}

	ok, err = s.repo.Users.Exists(chatIdAuthor)
	if err != nil {
		return 0, "", err
	}
	if !ok {
		return 0, "", PromoError
	}

	if chatIdAuthor == chatID {
		return 0, "", PromoError
	}

	err = s.UpdateSubscriptionDays(chatID, 7)
	if err != nil {
		return 0, "", err
	}

	err = s.UpdateSubscriptionDays(chatIdAuthor, 7)
	if err != nil {
		return 0, "", err
	}

	return chatIdAuthor, s.promo.Start, nil
}
