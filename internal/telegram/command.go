package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	log.WithFields(log.Fields{
		"chatID":  message.Chat.ID,
		"command": message.Command(),
	}).Info("new command")

	switch message.Command() {
	case b.cfg.Commands.Start:
		return b.handleStartCommand(message.Chat.ID, message.From.UserName, message.Text)
	case b.cfg.Commands.Help:
		return b.handleHelpCommand(message.Chat.ID)
	case b.cfg.Commands.Admin:
		return b.handleAdminCommand(message.Chat.ID)
	case b.cfg.Commands.Logs:
		return b.handleLogsCommand(message.Chat.ID)
	case b.cfg.Commands.DbUsers:
		return b.handleDbUsersCommand(message.Chat.ID)
	case b.cfg.Commands.DbUrls:
		return b.handleDbUrlsCommand(message.Chat.ID)
	case b.cfg.Commands.CreateUser:
		return b.handleCreateUserCommand(message.Chat.ID, message.Text)
	case b.cfg.Commands.CreateUrl:
		return b.handleCreateUrlCommand(message.Chat.ID, message.Text)
	case b.cfg.Commands.DeleteUser:
		return b.handleDeleteUserCommand(message.Chat.ID, message.Text)
	case b.cfg.Commands.DeleteUrl:
		return b.handleDeleteUrlCommand(message.Chat.ID, message.Text)
	case b.cfg.Commands.SendAll:
		return b.handleSendAllCommand(message.Chat.ID, message.Text)
	case b.cfg.Commands.SendAllBuy:
		return b.handleSendAllBuyCommand(message.Chat.ID, message.Text)
	case b.cfg.Commands.SendAllZeros:
		return b.handleSendAllZerosCommand(message.Chat.ID, message.Text)
	case b.cfg.Commands.GetSales:
		return b.handleGetSalesCommand(message.Chat.ID)
	case b.cfg.Commands.SetSales:
		return b.handleSetSalesCommand(message.Chat.ID, message.Text)
	case b.cfg.Commands.GetTokens:
		return b.handleGetTokensCommand(message.Chat.ID)
	case b.cfg.Commands.AddToken:
		return b.handleAddTokenCommand(message.Chat.ID, message.Text)
	case b.cfg.Commands.RemoveToken:
		return b.handleRemoveTokenCommand(message.Chat.ID, message.Text)
	case b.cfg.Commands.NextToken:
		return b.handleNextTokenCommand(message.Chat.ID)
	default:
		return invalidCommandError
	}
}

func (b *Bot) handleStartCommand(chatID int64, userName string, text string) error {
	text = strings.Replace(text, fmt.Sprintf("/%s", b.cfg.Commands.Start), "", 1)
	text = strings.Replace(text, " ", "", 1)
	if text == "" {
		text = "utm_empty"
	}

	err := b.services.CreateUser(chatID, userName, b.cfg.StartBalance, text)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.Start)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = b.NewMenuKeyboard()
	_, err = b.botApi.Send(msg)
	if err != nil {
		return err
	}

	msg = tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.AdviceButton)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = b.NewAdviceKeyboard()
	_, err = b.botApi.Send(msg)
	return err
}

func (b *Bot) handleHelpCommand(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.Help)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = b.NewAdviceKeyboard()
	_, err := b.botApi.Send(msg)
	return err
}
