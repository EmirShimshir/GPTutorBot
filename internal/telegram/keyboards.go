package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) NewMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(b.cfg.Keyboard.Menu.Help),
			tgbotapi.NewKeyboardButton(b.cfg.Keyboard.Menu.Balance),
			tgbotapi.NewKeyboardButton(b.cfg.Keyboard.Menu.Ref),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(b.cfg.Keyboard.Menu.ImageRecognize),
			tgbotapi.NewKeyboardButton(b.cfg.Keyboard.Menu.SolveTask),
		),
	)
}

func (b *Bot) NewShopKeyboard(chatID int64) tgbotapi.InlineKeyboardMarkup {
	product01, product02, product03 := b.services.GenerateProducts()
	url01, url02, url03 := b.services.GenerateURLs(chatID)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(product01, url01),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(product02, url02),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(product03, url03),
		),
	)
}

func (b *Bot) NewBalanceKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.cfg.Keyboard.Balance.Buy, b.cfg.Keyboard.Balance.Buy),
		),
			tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.cfg.Keyboard.Balance.Ref, b.cfg.Keyboard.Balance.Ref),
		),
	)
}

func (b *Bot) NewSalesKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔔 АКЦИЯ! 🔔", b.cfg.Keyboard.Balance.Buy),
		),
	)
}

func (b *Bot) NewRefKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.cfg.Keyboard.Balance.Ref, b.cfg.Keyboard.Balance.Ref),
		),
	)
}

func (b *Bot) NewAdviceKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.cfg.Keyboard.Advices.Advice, b.cfg.Keyboard.Advices.Advice),
		),
	)
}
