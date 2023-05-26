package service

import (
	"fmt"
	"github.com/EmirShimshir/tasker-bot/internal/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
)

func (s *Service) SetUrl(utm string, countRequests int64) error {
	log.WithFields(log.Fields{
		"utm": utm,
	}).Info("SetUrl")

	url := domain.NewUrl(utm)
	url.CountRequests = countRequests

	return s.repo.URLs.Save(url)
}

func (s *Service) getUrlsDataAll(urlsAll []*domain.Url, baseUrl string) ([]byte, error) {
	log.Info("getDataAll")

	file, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	defer os.Remove(file.Name())

	for _, u := range urlsAll {
		_, err = file.WriteString(fmt.Sprintf("%s;%d -> %s?start=%s\n", u.Utm, u.CountRequests, baseUrl, u.Utm))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	fileData, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (s *Service) GetAllUrls(baseUrl string) ([]byte, error) {
	log.Info("GetAllUrls")

	urlsAll, err := s.repo.URLs.GetAll()
	if err != nil {
		return nil, err
	}

	sort.Slice(urlsAll, func(i, j int) bool {
		return urlsAll[i].CountRequests < urlsAll[j].CountRequests
	})

	sortCountRequests, err := s.getUrlsDataAll(urlsAll, baseUrl)
	if err != nil {
		return nil, err
	}

	return sortCountRequests, nil
}

func (s *Service) DeleteUrl(utm string) error {
	log.Info("DeleteUrl")

	ok, err := s.repo.URLs.Exists(utm)
	if err != nil {
		return err
	}
	if !ok {
		return AdminUsageErr
	}

	return s.repo.URLs.Delete(utm)
}

func (s *Service) UpdateUtm(utm string) error {
	log.Info("UpdateUtm")

	ok, err := s.repo.URLs.Exists(utm)
	if err != nil {
		log.Error("err check exists utl", err)
		return err
	}
	if !ok {
		err := s.SetUrl(utm, 0)
		if err != nil {
			log.Error("err SetUrl", err)
			return err
		}
	}

	url, err := s.repo.URLs.Get(utm)
	if err != nil {
		log.Error("err s.repo.URLs.Get(utm)", err)
		return err
	}

	url.CountRequests++

	err = s.SetUrl(url.Utm, url.CountRequests)
	if err != nil {
		log.Error("err s.SetUrl(url.Utm, url.CountRequests)", err)
		return err
	}

	return nil
}
