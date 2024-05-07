package service

import (
	"fmt"
	"github.com/EmirShimshir/tasker-bot/internal/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *Service) CreateUser(chatID int64, userName string, startBalance int64, utm string) error {
	log.WithFields(log.Fields{
		"chatID": chatID,
		"Utm":    utm,
	}).Info("CreateUser")

	newUser := domain.NewUser(userName, chatID, startBalance)

	return s.repo.Users.Save(newUser)
}

func (s *Service) UserExists(chatID int64) (bool, error) {
	return s.repo.Users.Exists(chatID)
}

func (s *Service) GetUserBalance(chatID int64) (string, error) {
	log.WithFields(log.Fields{
		"chatID": chatID,
	}).Info("GetBalance")

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

	balance := strconv.FormatInt(user.Balance, 10)

	return balance, nil
}

func (s *Service) UpdateUserBalance(chatID int64, countBought int64) error {
	log.WithFields(log.Fields{
		"chatID":      chatID,
		"countBought": countBought,
	}).Info("UpdateBalance")

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

	user.Balance += countBought
	return s.repo.Users.Save(user)
}

func (s *Service) SetUser(chatID int64, newBalance int64, userSubDateEnd string, usedPromo bool) error {
	log.WithFields(log.Fields{
		"chatID":     chatID,
		"newBalance": newBalance,
	}).Info("SetBalance")

	ok, err := s.repo.Users.Exists(chatID)
	if err != nil {
		return err
	}
	if !ok {
		return AdminUsageErr
	}

	user, err := s.repo.Users.Get(chatID)
	if err != nil {
		return err
	}

	user.Balance = newBalance
	user.UsedPromo = usedPromo

	date := strings.Split(userSubDateEnd, ".")
	if len(date) != 3 {
		return AdminUsageErr
	}

	day, err := strconv.ParseInt(date[0], 10, 64)
	if err != nil {
		return err
	}

	month, err := strconv.ParseInt(date[1], 10, 64)
	if err != nil {
		return err
	}

	year, err := strconv.ParseInt(date[2], 10, 64)
	if err != nil {
		return err
	}

	user.DateSub = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Local)
	return s.repo.Users.Save(user)
}

func (s *Service) getUsersDataAll(usersAll []*domain.User, startMessage string) ([]byte, error) {
	log.Info("getDataAll")

	file, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	defer os.Remove(file.Name())

	_, err = file.WriteString(startMessage)
	if err != nil {
		return nil, err
	}

	for _, u := range usersAll {
		_, err = file.WriteString(fmt.Sprintf("@%s;%d;%d;%s;%t\n", u.Name, u.ChatID, u.Balance, u.DateSub.Format("02.01.2006"), u.UsedPromo))
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

func (s *Service) GetUsersAll(startMessage string) ([]byte, []byte, error) {
	log.Info("GetAllBalance")

	usersAll, err := s.repo.Users.GetAll()
	if err != nil {
		return nil, nil, err
	}

	sort.Slice(usersAll, func(i, j int) bool {
		return usersAll[i].Balance < usersAll[j].Balance
	})

	sortBalance, err := s.getUsersDataAll(usersAll, startMessage)
	if err != nil {
		return nil, nil, err
	}

	sort.Slice(usersAll, func(i, j int) bool {
		return usersAll[i].DateSub.Sub(usersAll[j].DateSub) < 0
	})

	sortDate, err := s.getUsersDataAll(usersAll, startMessage)
	if err != nil {
		return nil, nil, err
	}

	return sortBalance, sortDate, nil
}

func (s *Service) GetUsersAllChatID() ([]int64, error) {
	log.Info("GetAllChatID")

	usersAll, err := s.repo.Users.GetAll()
	if err != nil {
		return nil, err
	}

	ChatIDs := make([]int64, len(usersAll))

	for i := 0; i < len(ChatIDs); i++ {
		ChatIDs[i] = usersAll[i].ChatID
	}

	return ChatIDs, nil
}

func (s *Service) GetUsersZeroChatID() ([]int64, error) {
	log.Info("GetUsersZeroChatID")

	usersAll, err := s.repo.Users.GetAll()
	if err != nil {
		return nil, err
	}

	ChatIDs := make([]int64, 0, 1)

	now := time.Now()

	for i := 0; i < len(usersAll); i++ {
		diff := now.Sub(usersAll[i].DateSub)
		if diff > 0 && usersAll[i].Balance == 0 {
			if usersAll[i].UsedPromo == true {
				usersAll[i].UsedPromo = false
				if err := s.repo.Users.Save(usersAll[i]); err != nil {
					return nil, err
				}
			}
			ChatIDs = append(ChatIDs, usersAll[i].ChatID)
		}
	}

	return ChatIDs, nil
}

func (s *Service) GetUsersNotSubChatID() ([]int64, error) {
	log.Info("GetUsersNotSubChatID")

	usersAll, err := s.repo.Users.GetAll()
	if err != nil {
		return nil, err
	}

	ChatIDs := make([]int64, 0, 1)

	now := time.Now()

	for i := 0; i < len(usersAll); i++ {
		diff := now.Sub(usersAll[i].DateSub)
		if diff > 0 {
			ChatIDs = append(ChatIDs, usersAll[i].ChatID)
		}
	}

	return ChatIDs, nil
}

func (s *Service) DeleteUser(chatID int64) error {
	log.Info("DeleteUser")

	ok, err := s.repo.Users.Exists(chatID)
	if err != nil {
		return err
	}
	if !ok {
		return AdminUsageErr
	}

	return s.repo.Users.Delete(chatID)
}
