package telegram

import (
	"errors"
	"fmt"
	"github.com/EmirShimshir/tasker-bot/internal/chatAI"
	"github.com/EmirShimshir/tasker-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"time"
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
	case b.cfg.Keyboard.Menu.Ref:
		return b.handleCreateRef(chatID)
	case b.cfg.Keyboard.Menu.ImageRecognize:
		return b.handleImage(chatID)
	case b.cfg.Keyboard.Menu.SolveTask:
		return b.handleTask(chatID)
	}

	if b.services.IsGift(text) {
		return b.handleActivateGift(chatID)
	}

	return b.handleRawText(text, chatID)
}

func (b *Bot) typeWhileChannelOpen(stopCh <-chan struct{}, chatID int64) {
	for {
		select {
		case <-stopCh:
			return
		default:
			action := tgbotapi.ChatActionConfig{BaseChat: tgbotapi.BaseChat{ChatID: chatID}, Action: "typing"}
			_, _ = b.botApi.Request(action)
			time.Sleep(5 * time.Second)
		}

	}
}

func (b *Bot) handleRawText(text string, chatID int64) error {
	stopCh := make(chan struct{})

	go b.typeWhileChannelOpen(stopCh, chatID)

	result, err := b.services.ProcessTask(text, chatID)
	stopCh <- struct{}{}

	if errors.Is(err, chatAI.ErrGptResult) || errors.Is(err, chatAI.GptNewToken) {
		AdminText := ""
		if errors.Is(err, chatAI.ErrGptResult) {
			AdminText = fmt.Sprintf("ALERT: %s", err.Error())
		} else {
			AdminText = fmt.Sprintf("INFO: %s", err.Error())
		}
		for _, adminID := range b.cfg.AdminsId {
			msg := tgbotapi.NewMessage(adminID, AdminText)
			_, err = b.botApi.Send(msg)
		}
	} else if errors.Is(err, service.EmptyBalanceError) {
		return b.handleBalance(chatID)
	} else if err != nil {
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
