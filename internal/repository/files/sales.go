package files

import (
	"encoding/gob"
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"os"
	"path/filepath"
)

type SalesRepo struct {
	basePath string
}

const (
	SalesDefaultPerm = 0774
)

func NewSalesRepo(repo config.Repo) *SalesRepo {
	return &SalesRepo{
		basePath: repo.SalesBasePath,
	}
}

func (s *SalesRepo) Save(count int64) error {
	err := os.MkdirAll(s.basePath, UsersDefaultPerm)
	if err != nil {
		return err
	}

	filePath := filepath.Join(s.basePath, "count")

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err = gob.NewEncoder(file).Encode(count); err != nil {
		return err
	}

	return nil
}

func (s *SalesRepo) Get() (int64, error) {
	filePath := filepath.Join(s.basePath, "count")

	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer func() { _ = file.Close() }()

	count := int64(0)

	err = gob.NewDecoder(file).Decode(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
