package telegram

import (
	"fmt"
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"github.com/EmirShimshir/tasker-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
)

type Bot struct {
	botApi   *tgbotapi.BotAPI
	services *service.Service
	cfg      config.Bot
}

func NewBot(botApi *tgbotapi.BotAPI, services *service.Service, cfg config.Bot) *Bot {
	return &Bot{
		botApi:   botApi,
		services: services,
		cfg:      cfg,
	}
}

func (b *Bot) Start() error {
	wh, err := tgbotapi.NewWebhookWithCert(fmt.Sprintf("https://%s:%d/%s", b.cfg.Ip, b.cfg.Port, b.botApi.Token), tgbotapi.FilePath(b.cfg.CertPath))
	if err != nil {
		return err
	}

	_, err = b.botApi.Request(wh)
	if err != nil {
		return err
	}

	updates := b.botApi.ListenForWebhook(fmt.Sprintf("/%s", b.botApi.Token))
	go http.ListenAndServeTLS(fmt.Sprintf(":%d", b.cfg.Port), b.cfg.CertPath, b.cfg.KeyPath, nil)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				err := b.handleCommand(update.Message)
				if err != nil {
					b.handleError(update.Message.Chat.ID, err)
				}
			} else {
				err := b.handleMessage(update.Message)
				if err != nil {
					b.handleError(update.Message.Chat.ID, err)
				}
			}
		} else if update.CallbackQuery != nil {
			err := b.handleCallback(update.CallbackQuery)
			if err != nil {
				b.handleError(update.CallbackQuery.Message.Chat.ID, err)
			}
		}
	}

	return nil
}
