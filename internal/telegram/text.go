package telegram

import (
	"errors"
	"fmt"
	"github.com/EmirShimshir/tasker-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handleText(text string, chatID int64) error {
	log.WithFields(log.Fields{
		"chatID": chatID,
		"text":   text,
	}).Info("new text")

	switch text {
	case b.cfg.Keyboard.Menu.Help:
		return b.handleHelpCommand(chatID)
	case b.cfg.Keyboard.Menu.Balance:
		return b.handleBalance(chatID)
	case b.cfg.Keyboard.Menu.Promo:
		return b.handlePromo(chatID)
	case b.cfg.Keyboard.Menu.ImageRecognize:
		return b.handleImage(chatID)
	case b.cfg.Keyboard.Menu.SolveTask:
		return b.handleTask(chatID)
	}

	if b.services.IsPromo(text) {
		return b.handleActivatePromo(text, chatID)
	}

	if b.services.IsGift(text) {
		return b.handleActivateGift(chatID)
	}

	return b.handleRawText(text, chatID)
}

func (b *Bot) handleRawText(text string, chatID int64) error {
	action := tgbotapi.ChatActionConfig{tgbotapi.BaseChat{ChatID: chatID}, "typing"}
	_, err := b.botApi.Request(action)
	if err != nil {
		return err
	}

	result, err := b.services.ProcessTask(text, chatID)
	if errors.Is(err, service.EmptyBalanceError) {
		return b.handleBalance(chatID)
	}
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(chatID, result)

	_, err = b.botApi.Send(msg)
	if err != nil {
		return err
	}
	return b.handleBalance(chatID)
}

func (b *Bot) handleBalance(chatID int64) error {
	ok, err := b.services.IsSubscriber(chatID)
	if err != nil {
		return err
	}

	text := ""
	if ok {
		dateEnd, err := b.services.GetSubscribeDateEnd(chatID)
		if err != nil {
			return err
		}
		text = fmt.Sprintf("%s*%s*", b.cfg.Messages.Responses.Subscribe, dateEnd)
	} else {
		balance, err := b.services.GetUserBalance(chatID)
		if err != nil {
			return err
		}
		text = fmt.Sprintf("%s*%s*", b.cfg.Messages.Responses.Balance, balance)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = b.NewBalanceKeyboard()
	_, err = b.botApi.Send(msg)
	return err
}

func (b *Bot) handleImage(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.Image)
	_, err := b.botApi.Send(msg)
	return err
}

func (b *Bot) handleTask(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.Task)
	_, err := b.botApi.Send(msg)
	return err
}

func (b *Bot) handlePromo(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Promo)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = b.NewPromoKeyboard()
	_, err := b.botApi.Send(msg)
	if err != nil {
		return err
	}

	return err
}
