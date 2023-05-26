package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	if message.Document != nil {
		fileConfig := tgbotapi.FileConfig{FileID: message.Document.FileID}
		return b.handleFile(fileConfig, message.Chat.ID)
	}
	if message.Photo != nil {
		photo := message.Photo[len(message.Photo)-1]
		fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
		return b.handleFile(fileConfig, message.Chat.ID)
	}
	if message.Text != "" {
		return b.handleText(message.Text, message.Chat.ID)
	}

	return invalidMessageError
}

func (b *Bot) handleFile(fileConfig tgbotapi.FileConfig, chatID int64) error {
	log.WithFields(log.Fields{
		"chatID": chatID,
	}).Info("new file")

	action := tgbotapi.ChatActionConfig{tgbotapi.BaseChat{ChatID: chatID}, "typing"}
	_, err := b.botApi.Request(action)
	if err != nil {
		return err
	}

	tgFile, err := b.botApi.GetFile(fileConfig)
	if err != nil {
		return err
	}

	urlFile := tgFile.Link(b.botApi.Token)

	text, err := b.services.ProcessFile(urlFile)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, text)
	_, err = b.botApi.Send(msg)
	return err
}
