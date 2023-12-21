package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
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

func (b *Bot) handleSendAllBuyCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	text = strings.Replace(text, fmt.Sprintf("/%s ", b.cfg.Commands.SendAllBuy), "", 1)

	// send to admin
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = b.NewSalesKeyboard()
	_, err := b.botApi.Send(msg)
	if err != nil {
		log.Error(fmt.Sprintf("err sendAll to id: %d", chatID), err)
	}
	log.WithFields(log.Fields{
		"chatID": chatID,
		"status": "OK",
	}).Info("handleSendAllCommand to admin")

	chatIDs, err := b.services.GetUsersAllChatID()
	if err != nil {
		return err
	}

	// send all
	for _, chatID := range chatIDs {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = b.NewSalesKeyboard()
		_, err = b.botApi.Send(msg)
		if err != nil {
			log.Error(fmt.Sprintf("err sendAll to id: %d", chatID), err)
		}
		log.WithFields(log.Fields{
			"chatID": chatID,
			"status": "OK",
		}).Info("handleSendAllCommand")
	}

	log.WithFields(log.Fields{
		"status": "done",
	}).Info("handleSendAllCommand")
	return nil
}

func (b *Bot) handleSendAllZerosCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	text = strings.Replace(text, fmt.Sprintf("/%s ", b.cfg.Commands.SendAllZeros), "", 1)

	// send to admin
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := b.botApi.Send(msg)
	if err != nil {
		log.Error(fmt.Sprintf("err SendAllZeros to id: %d", chatID), err)
	}
	log.WithFields(log.Fields{
		"chatID": chatID,
		"status": "OK",
	}).Info("handleSendAllZerosCommand to admin")

	chatIDs, err := b.services.GetUsersZeroChatID()
	if err != nil {
		return err
	}

	// send all
	for _, chatID := range chatIDs {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		_, err = b.botApi.Send(msg)
		if err != nil {
			log.Error(fmt.Sprintf("err SendAllZeros to id: %d", chatID), err)
		}
		log.WithFields(log.Fields{
			"chatID": chatID,
			"status": "OK",
		}).Info("handleSendAllZerosCommand")
	}

	log.WithFields(log.Fields{
		"status": "done",
	}).Info("handleSendAllZerosCommand")
	return nil
}

func (b *Bot) handleSendAllCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	text = strings.Replace(text, fmt.Sprintf("/%s ", b.cfg.Commands.SendAll), "", 1)

	// send to admin
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := b.botApi.Send(msg)
	if err != nil {
		log.Error(fmt.Sprintf("err sendAll to id: %d", chatID), err)
	}
	log.WithFields(log.Fields{
		"chatID": chatID,
		"status": "OK",
	}).Info("handleSendAllCommand to admin")

	chatIDs, err := b.services.GetUsersAllChatID()
	if err != nil {
		return err
	}

	// send all
	for _, chatID := range chatIDs {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		_, err = b.botApi.Send(msg)
		if err != nil {
			log.Error(fmt.Sprintf("err sendAll to id: %d", chatID), err)
		}
		log.WithFields(log.Fields{
			"chatID": chatID,
			"status": "OK",
		}).Info("handleSendAllCommand")
	}

	log.WithFields(log.Fields{
		"status": "done",
	}).Info("handleSendAllCommand")
	return nil
}

func (b *Bot) handleGetSalesCommand(chatID int64) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	count := b.services.GetSales()

	text := fmt.Sprintf("count: %d", count)

	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.botApi.Send(msg)

	log.WithFields(log.Fields{
		"msg":    text,
		"status": "done",
	}).Info("handleGetSalesCommand")
	return err
}

func (b *Bot) handleSetSalesCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	text = strings.Replace(text, fmt.Sprintf("/%s ", b.cfg.Commands.SetSales), "", 1)

	count, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return err
	}

	b.services.SetSales(count)

	text = fmt.Sprintf("seted: %d", count)

	msg := tgbotapi.NewMessage(chatID, text)
	_, err = b.botApi.Send(msg)

	log.WithFields(log.Fields{
		"msg":    count,
		"status": "done",
	}).Info("handleSetSalesCommand")
	return err
}

func (b *Bot) handleGetTokensCommand(chatID int64) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	data, err := b.services.GetTokensDataAll()
	if err != nil {
		return err
	}

	document := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{"tokens.txt", data})
	_, err = b.botApi.Send(document)
	return err
}

func (b *Bot) handleAddTokenCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	text = strings.Replace(text, fmt.Sprintf("/%s ", b.cfg.Commands.AddToken), "", 1)

	b.services.AddToken(text)

	return nil
}

func (b *Bot) handleRemoveTokenCommand(chatID int64, text string) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	text = strings.Replace(text, fmt.Sprintf("/%s ", b.cfg.Commands.RemoveToken), "", 1)

	return	b.services.RemoveToken(text)
}

func (b *Bot) handleNextTokenCommand(chatID int64) error {
	if !b.isAdmin(chatID) {
		return invalidCommandError
	}

	b.services.NextToken()

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
