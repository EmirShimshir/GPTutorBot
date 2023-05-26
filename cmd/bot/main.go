package main

import (
	"fmt"
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"github.com/EmirShimshir/tasker-bot/internal/openai"
	"github.com/EmirShimshir/tasker-bot/internal/repository"
	"github.com/EmirShimshir/tasker-bot/internal/server"
	"github.com/EmirShimshir/tasker-bot/internal/service"
	"github.com/EmirShimshir/tasker-bot/internal/telegram"
	"github.com/EmirShimshir/tasker-bot/internal/tesseract"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	defaultPerm = 0777
	logsDir     = "./logs"
	logsPath    = "./logs/logs.txt"
)

var f *os.File

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	err := os.MkdirAll(logsDir, defaultPerm)
	if err != nil {
		fmt.Println("Failed to create log dir")
		panic(err)
	}

	f, err := os.OpenFile(logsPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, defaultPerm)
	if err != nil {
		fmt.Println("Failed to create log file")
		panic(err)
	}

	log.SetOutput(f)

	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("application startup...")

	defer f.Close()
	log.Info("logger initialized")

	cfg, err := config.Init()
	if err != nil {
		log.WithFields(log.Fields{
			"from":    "main()",
			"problem": "can't initialize config",
		}).Fatal(err.Error())
	}
	log.Info("config initialized")

	repo := repository.NewRepositories(cfg.Repo)
	log.Info("repositories initialized")

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.WithFields(log.Fields{
			"from":    "main()",
			"problem": "can't initialize botApi",
		}).Fatal(err.Error())
	}
	log.Info("botApi initialized")

	nlp := tesseract.NewNlp(cfg.Tesseract)
	defer nlp.CloseClient()
	log.Info("nlp initialized")

	chatGpt := openai.NewChatGpt(cfg.ChatGptApiKey, cfg.ChatGpt)
	log.Info("chatGpt initialized")

	services := service.NewService(nlp, chatGpt, cfg.Payment, cfg.Promo, repo)
	log.Info("services initialized")

	bot := telegram.NewBot(botApi, services, cfg.Bot)
	log.Info("bot initialized")

	serverPayment := server.NewPayment(services, cfg.Server)
	log.Info("server payment initialized")

	log.Info("server payment started")
	go func() {
		if err := serverPayment.Start(); err != nil {
			log.WithFields(log.Fields{
				"from":    "main()",
				"problem": "server payment shutdown",
			}).Fatal(err.Error())
		}
	}()

	log.Info("bot started")
	if err := bot.Start(); err != nil {
		log.WithFields(log.Fields{
			"from":    "main()",
			"problem": "bot shutdown",
		}).Fatal(err.Error())
	}
}
