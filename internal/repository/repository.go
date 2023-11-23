package repository

import (
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"github.com/EmirShimshir/tasker-bot/internal/domain"
	"github.com/EmirShimshir/tasker-bot/internal/repository/files"
)

type Users interface {
	Save(u *domain.User) error
	Get(chatID int64) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Exists(chatID int64) (bool, error)
	Delete(chatID int64) error
}

type URLs interface {
	Save(u *domain.Url) error
	Get(utm string) (*domain.Url, error)
	GetAll() ([]*domain.Url, error)
	Exists(utm string) (bool, error)
	Delete(utm string) error
}

type Sales interface {
	Save(count int64) error
	Get() (int64, error)
}
type Repositories struct {
	Users Users
	URLs  URLs
	Sales Sales
}

func NewRepositories(repo config.Repo) *Repositories {
	return &Repositories{
		Users: files.NewUsersRepo(repo),
		URLs:  files.NewUrlsRepo(repo),
		Sales: files.NewSalesRepo(repo),
	}
}
