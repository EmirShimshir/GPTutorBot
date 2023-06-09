package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"strconv"
	"strings"
)

func (b *Bot) handleAdminCommand(chatID int64) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.Admin)
	_, err := b.botApi.Send(msg)
	return err
}

func (b *Bot) handleLogsCommand(chatID int64) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	fileData, err := ioutil.ReadFile("./logs/logs.txt")
	if err != nil {
		return err
	}

	document := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{"logs.txt", fileData})
	_, err = b.botApi.Send(document)
	return err
}

func (b *Bot) handleDbUsersCommand(chatID int64) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	sortBalance, sortDate, err := b.services.GetUsersAll(b.cfg.Messages.Responses.AdminDbUsers)
	if err != nil {
		return err
	}

	document := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{"users_sort_balance.txt", sortBalance})
	_, err = b.botApi.Send(document)
	if err != nil {
		return err
	}

	document = tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{"users_sort_date.txt", sortDate})
	_, err = b.botApi.Send(document)
	return err
}

func (b *Bot) handleDbUrlsCommand(chatID int64) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	sortCountRequests, err := b.services.GetAllUrls(b.cfg.Url)
	if err != nil {
		return err
	}

	document := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{"urls_sort_count_requests.txt", sortCountRequests})
	_, err = b.botApi.Send(document)
	return err
}

func (b *Bot) handleCreateUserCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	arrParams := strings.Split(text, ";")
	if len(arrParams) != 5 {
		return invalidCommandError
	}

	userChatID, err := strconv.ParseInt(arrParams[1], 10, 64)
	if err != nil {
		return invalidCommandError
	}
	userBalance, err := strconv.ParseInt(arrParams[2], 10, 64)
	if err != nil {
		return invalidCommandError
	}
	userSubDateEnd := arrParams[3]
	if err != nil {
		return invalidCommandError
	}
	usedPromo, err := strconv.ParseBool(arrParams[4])
	if err != nil {
		return invalidCommandError
	}

	return b.services.SetUser(userChatID, userBalance, userSubDateEnd, usedPromo)
}

func (b *Bot) handleCreateUrlCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	text = strings.Replace(text, fmt.Sprintf("/%s ", b.cfg.Commands.CreateUrl), "", 1)
	arrParams := strings.Split(text, ";")
	if len(arrParams) != 2 {
		return invalidCommandError
	}

	utm := arrParams[0]

	countRequests, err := strconv.ParseInt(arrParams[1], 10, 64)
	if err != nil {
		return invalidCommandError
	}

	err = b.services.SetUrl(utm, countRequests)
	if err != nil {
		return invalidCommandError
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s?start=%s", b.cfg.Url, utm))
	_, err = b.botApi.Send(msg)
	return err
}

func (b *Bot) handleDeleteUserCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	text = strings.Replace(text, fmt.Sprintf("/%s ", b.cfg.Commands.DeleteUser), "", 1)

	userChatID, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return invalidCommandError
	}

	return b.services.DeleteUser(userChatID)
}

func (b *Bot) handleDeleteUrlCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	text = strings.Replace(text, fmt.Sprintf("/%s ", b.cfg.Commands.DeleteUrl), "", 1)

	utm := text

	return b.services.DeleteUrl(utm)
}

func (b *Bot) handleSendAllCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	chatIDs, err := b.services.GetUsersAllChatID()
	if err != nil {
		return err
	}

	text = strings.Replace(text, fmt.Sprintf("/%s ", b.cfg.Commands.SendAll), "", 1)
	for _, chatID := range chatIDs {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		_, _ = b.botApi.Send(msg)
	}

	return nil
}

func (b *Bot) isAdmin(chatID int64) bool {
	for _, adminId := range b.cfg.AdminsId {
		if chatID == adminId {
			return true
		}
	}
	return false
}
