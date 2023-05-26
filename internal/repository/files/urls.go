package files

import (
	"encoding/gob"
	"errors"
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"github.com/EmirShimshir/tasker-bot/internal/domain"
	"os"
	"path/filepath"
)

type UrlsRepo struct {
	basePath string
}

const (
	UrlsDefaultPerm = 0774
)

func NewUrlsRepo(repo config.Repo) *UrlsRepo {
	return &UrlsRepo{
		basePath: repo.UrlsBasePath,
	}
}

func (u *UrlsRepo) Save(url *domain.Url) error {
	err := os.MkdirAll(u.basePath, UrlsDefaultPerm)
	if err != nil {
		return err
	}

	filePath := filepath.Join(u.basePath, url.Utm)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err = gob.NewEncoder(file).Encode(url); err != nil {
		return err
	}

	return nil
}

func (u *UrlsRepo) Get(utm string) (*domain.Url, error) {
	filePath := filepath.Join(u.basePath, utm)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	url := new(domain.Url)

	err = gob.NewDecoder(file).Decode(url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (u *UrlsRepo) GetAll() ([]*domain.Url, error) {
	UrlsAll := make([]*domain.Url, 0, 1)

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

		url := new(domain.Url)

		err = gob.NewDecoder(file).Decode(url)
		if err != nil {
			return nil, err
		}

		UrlsAll = append(UrlsAll, url)
	}

	return UrlsAll, nil
}

func (u *UrlsRepo) Exists(utm string) (bool, error) {
	filePath := filepath.Join(u.basePath, utm)

	switch _, err := os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, err
	}

	return true, nil
}

func (u *UrlsRepo) Delete(utm string) error {
	filePath := filepath.Join(u.basePath, utm)

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}
