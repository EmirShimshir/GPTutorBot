package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) NewMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(b.cfg.Keyboard.Menu.Help),
			tgbotapi.NewKeyboardButton(b.cfg.Keyboard.Menu.Balance),
			tgbotapi.NewKeyboardButton("Промокод"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(b.cfg.Keyboard.Menu.ImageRecognize),
			tgbotapi.NewKeyboardButton(b.cfg.Keyboard.Menu.SolveTask),
		),
	)
}

func (b *Bot) NewShopKeyboard(chatID int64) tgbotapi.InlineKeyboardMarkup {
	product01 := b.services.GenerateProduct(b.cfg.Shop.ProductCount01, b.cfg.Shop.ProductPrice01)
	product02 := b.services.GenerateProduct(b.cfg.Shop.ProductCount02, b.cfg.Shop.ProductPrice02)
	product03 := b.services.GenerateProduct(b.cfg.Shop.ProductCount03, b.cfg.Shop.ProductPrice03)

	url01 := b.services.GenerateURL(chatID, b.cfg.Shop.ProductCount01, b.cfg.Shop.ProductPrice01)
	url02 := b.services.GenerateURL(chatID, b.cfg.Shop.ProductCount02, b.cfg.Shop.ProductPrice02)
	url03 := b.services.GenerateURL(chatID, b.cfg.Shop.ProductCount03, b.cfg.Shop.ProductPrice03)

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
	)
}

func (b *Bot) NewAdviceKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.cfg.Keyboard.Advices.Advice, b.cfg.Keyboard.Advices.Advice),
		),
	)
}

func (b *Bot) NewPromoKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(b.cfg.Keyboard.Promo.CreatePromo, b.cfg.Keyboard.Promo.CreatePromo),
		),
	)
}
