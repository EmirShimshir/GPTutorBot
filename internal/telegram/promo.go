package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCreatePromo(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.CreatePromo01)
	msg.ParseMode = "Markdown"
	_, err := b.botApi.Send(msg)
	if err != nil {
		return err
	}

	promo, err := b.services.CreatePromo(chatID)
	if err != nil {
		return err
	}

	msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("*%s*", promo))
	msg.ParseMode = "Markdown"
	_, err = b.botApi.Send(msg)
	if err != nil {
		return err
	}

	msg = tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.CreatePromo02)
	msg.ParseMode = "Markdown"
	_, err = b.botApi.Send(msg)
	if err != nil {
		return err
	}

	msg = tgbotapi.NewMessage(chatID, fmt.Sprintf(b.cfg.Messages.Responses.CreatePromo03, promo))
	msg.ParseMode = "Markdown"
	_, err = b.botApi.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleActivatePromo(text string, chatID int64) error {
	err := b.services.ActivatePromo(text, chatID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.PromoOK)
	msg.ParseMode = "Markdown"
	_, err = b.botApi.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
