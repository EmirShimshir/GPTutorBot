package telegram

import (
	"errors"
	"github.com/EmirShimshir/tasker-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

var (
	invalidMessageError  = errors.New("message is invalid")
	invalidCommandError  = errors.New("command is invalid")
	InvalidCallbackError = errors.New("callback is invalid")
	AdminErr             = errors.New("it is not admin")
)

func (b *Bot) handleError(chatID int64, err error) {
	var messageText string

	log.WithFields(log.Fields{
		"chatID": chatID,
	}).Error(err)

	switch err {
	case invalidMessageError:
		messageText = b.cfg.Messages.Errors.UnknownMessage
	case invalidCommandError:
		messageText = b.cfg.Messages.Errors.UnknownCommand
	case InvalidCallbackError:
		messageText = b.cfg.Messages.Errors.UnknownCallback
	case AdminErr:
		messageText = b.cfg.Messages.Errors.AdminError
	case service.AdminUsageErr:
		messageText = b.cfg.Messages.Errors.AdminUsageError
	case service.NotAuthError:
		messageText = b.cfg.Messages.Errors.NotAuth
	case service.PromoError:
		messageText = b.cfg.Messages.Errors.PromoError
	case service.PromoUsedError:
		messageText = b.cfg.Messages.Errors.PromoUsedError
	default:
		messageText = b.cfg.Messages.Errors.Default
	}

	msg := tgbotapi.NewMessage(chatID, messageText)
	_, _ = b.botApi.Send(msg)
}
