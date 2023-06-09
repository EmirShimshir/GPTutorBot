package service

import (
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"github.com/EmirShimshir/tasker-bot/internal/openai"
	"github.com/EmirShimshir/tasker-bot/internal/repository"
	"github.com/EmirShimshir/tasker-bot/internal/tesseract"
)

type Service struct {
	nlp     *tesseract.Nlp
	chatGpt *openai.ChatGpt
	payment config.Payment
	promo   config.Promo
	repo    *repository.Repositories
}

func NewService(nlp *tesseract.Nlp, chatGpt *openai.ChatGpt, payment config.Payment, promo config.Promo, repo *repository.Repositories) *Service {
	return &Service{
		nlp:     nlp,
		chatGpt: chatGpt,
		payment: payment,
		promo:   promo,
		repo:    repo,
	}
}
