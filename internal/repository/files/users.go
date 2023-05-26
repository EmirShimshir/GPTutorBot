package files

import (
	"encoding/gob"
	"errors"
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"github.com/EmirShimshir/tasker-bot/internal/domain"
	"os"
	"path/filepath"
	"strconv"
)

type UsersRepo struct {
	basePath string
}

const (
	UsersDefaultPerm = 0774
)

func NewUsersRepo(repo config.Repo) *UsersRepo {
	return &UsersRepo{
		basePath: repo.UsersBasePath,
	}
}

func (u *UsersRepo) Save(user *domain.User) error {
	err := os.MkdirAll(u.basePath, UsersDefaultPerm)
	if err != nil {
		return err
	}

	fileName := strconv.FormatInt(user.ChatID, 10)
	filePath := filepath.Join(u.basePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err = gob.NewEncoder(file).Encode(user); err != nil {
		return err
	}

	return nil
}

func (u *UsersRepo) Get(chatID int64) (*domain.User, error) {
	fileName := strconv.FormatInt(chatID, 10)
	filePath := filepath.Join(u.basePath, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	user := new(domain.User)

	err = gob.NewDecoder(file).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UsersRepo) GetAll() ([]*domain.User, error) {
	usersAll := make([]*domain.User, 0, 1)

	files, err := os.ReadDir(u.basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(u.basePath, fileName)

		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer func() { _ = file.Close() }()

		user := new(domain.User)

		err = gob.NewDecoder(file).Decode(user)
		if err != nil {
			return nil, err
		}

		usersAll = append(usersAll, user)
	}

	return usersAll, nil
}

func (u *UsersRepo) Exists(chatID int64) (bool, error) {
	fileName := strconv.FormatInt(chatID, 10)
	filePath := filepath.Join(u.basePath, fileName)

	switch _, err := os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, err
	}

	return true, nil
}

func (u *UsersRepo) Delete(chatID int64) error {
	fileName := strconv.FormatInt(chatID, 10)
	filePath := filepath.Join(u.basePath, fileName)

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}
