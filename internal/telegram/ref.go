package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCreateRef(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.CreateRef01)
	msg.ParseMode = "Markdown"
	_, err := b.botApi.Send(msg)
	if err != nil {
		return err
	}

	ref, err := b.services.CreateRef(chatID)
	if err != nil {
		return err
	}

	msg = tgbotapi.NewMessage(chatID, ref)
	msg.ParseMode = "Markdown"
	_, err = b.botApi.Send(msg)
	if err != nil {
		return err
	}

	msg = tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.CreateRef02)
	msg.ParseMode = "Markdown"
	_, err = b.botApi.Send(msg)
	if err != nil {
		return err
	}

	msg = tgbotapi.NewMessage(chatID, fmt.Sprintf(b.cfg.Messages.Responses.CreateRef03, ref))
	msg.ParseMode = "Markdown"
	_, err = b.botApi.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleActivateRef(text string, chatID int64) (string, error) {
	chatIdAuthor, refStart, err := b.services.ActivateRef(text, chatID)
	if err != nil {
		return "", err
	}

	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.PromoOK)
	msg.ParseMode = "Markdown"
	_, err = b.botApi.Send(msg)
	if err != nil {
		return "", err
	}

	msg = tgbotapi.NewMessage(chatIdAuthor, b.cfg.Messages.Responses.PromoOK)
	msg.ParseMode = "Markdown"
	_, err = b.botApi.Send(msg)
	if err != nil {
		return "", err
	}

	return refStart, nil
}
