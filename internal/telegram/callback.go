package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) error {
	log.WithFields(log.Fields{
		"chatID": callback.Message.Chat.ID,
		"data":   callback.Data,
	}).Info("new callback")

	switch callback.Data {
	case b.cfg.Keyboard.Balance.Buy:
		return b.handleBuyRequests(callback.Message.Chat.ID)
	case b.cfg.Keyboard.Advices.Advice:
		return b.handleExamples(callback.Message.Chat.ID)
	case b.cfg.Keyboard.Promo.CreatePromo:
		return b.handleCreatePromo(callback.Message.Chat.ID)

	default:
		return InvalidCallbackError
	}
}

func (b *Bot) handleBuyRequests(chatID int64) error {
	count, err := b.services.GetSales()
	if err != nil {
		return err
	}

	var text string

	if count > 0 {
		count, err := b.services.GetSales()
		if err != nil {
			return err
		}
		text = fmt.Sprintf(b.cfg.Messages.Responses.BuyRequestsSales, count)
		fmt.Println(text)

	} else {
		text = b.cfg.Messages.Responses.BuyRequests

	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = b.NewShopKeyboard(chatID)

	_, err = b.botApi.Send(msg)
	return err
}

func (b *Bot) handleExamples(chatID int64) error {
	adviceStart := tgbotapi.NewMessage(chatID, b.cfg.Messages.Responses.AdviceStart)
	adviceStart.ParseMode = "Markdown"
	_, err := b.botApi.Send(adviceStart)
	if err != nil {
		return err
	}

	advice01 := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("./images/advice_01.jpeg"))
	advice01.Caption = b.cfg.Messages.Responses.Advice01

	_, err = b.botApi.Send(advice01)
	if err != nil {
		return err
	}

	advice02 := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("./images/advice_02.jpeg"))
	advice02.Caption = b.cfg.Messages.Responses.Advice02

	_, err = b.botApi.Send(advice02)
	if err != nil {
		return err
	}

	advice03 := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("./images/advice_03.jpeg"))
	advice03.Caption = b.cfg.Messages.Responses.Advice03

	_, err = b.botApi.Send(advice03)
	if err != nil {
		return err
	}

	advice04 := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("./images/advice_04.jpeg"))
	advice04.Caption = b.cfg.Messages.Responses.Advice04

	_, err = b.botApi.Send(advice04)
	if err != nil {
		return err
	}

	return nil
}
