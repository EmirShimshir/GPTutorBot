package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (s *Service) CreatePromo(chatID int64) (string, error) {
	log.WithFields(log.Fields{
		"chatID": chatID,
	}).Info("CreatePromo")

	promo := fmt.Sprintf("%s%d", s.promo.Start, chatID)

	return promo, nil
}

func (s *Service) IsPromo(text string) bool {

	return strings.HasPrefix(text, s.promo.Start)
}

func (s *Service) ActivatePromo(text string, chatID int64) error {
	log.WithFields(log.Fields{
		"chatID": chatID,
	}).Info("ActivatePromo")

	argsArray := strings.Split(text, ":")
	if len(argsArray) != 2 {
		return PromoError
	}

	chatIdAuthor, err := strconv.ParseInt(argsArray[1], 10, 64)
	if err != nil {
		return PromoError
	}

	ok, err := s.repo.Users.Exists(chatID)
	if err != nil {
		return err
	}
	if !ok {
		return NotAuthError
	}

	ok, err = s.repo.Users.Exists(chatIdAuthor)
	if err != nil {
		return err
	}
	if !ok {
		return PromoError
	}

	if chatIdAuthor == chatID {
		return PromoError
	}

	err = s.UpdatePromoBalanceFriend(chatID, s.promo.CountReqFriend)
	if err != nil {
		return err
	}

	err = s.UpdatePromoBalanceAuthor(chatIdAuthor, s.promo.CountReqAuthor)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdatePromoBalanceAuthor(chatID int64, countReq int64) error {
	user, err := s.repo.Users.Get(chatID)
	if err != nil {
		return err
	}

	user.Balance += countReq

	return s.repo.Users.Save(user)
}

func (s *Service) UpdatePromoBalanceFriend(chatId int64, countReq int64) error {
	user, err := s.repo.Users.Get(chatId)
	if err != nil {
		return err
	}

	if user.UsedPromo == false {
		user.Balance += countReq
		user.UsedPromo = true
	} else {
		return PromoUsedError
	}

	return s.repo.Users.Save(user)
}
