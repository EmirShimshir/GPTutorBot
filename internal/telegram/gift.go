package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) handleActivateGift(chatID int64) error {
	err := b.services.ActivateGift(chatID)
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
