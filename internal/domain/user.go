package domain

import "time"

type User struct {
	Name      string
	ChatID    int64
	Balance   int64
	DateSub   time.Time
	UsedPromo bool
}

func NewUser(name string, chatID int64, balance int64) *User {
	return &User{
		Name:      name,
		ChatID:    chatID,
		Balance:   balance,
		DateSub:   time.Now(),
		UsedPromo: false,
	}
}
