package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"strings"
)

func (s *Service) GenerateURL(chatID int64, countBought int64, price int64) string {
	base := url.URL{
		Scheme: s.payment.URL.Scheme,
		Host:   s.payment.URL.Host,
		Path:   s.payment.URL.Path,
	}

	uEnc, err := s.encryptProductAES(chatID, countBought)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("encryptProductAES")
		return ""
	}

	params := url.Values{
		"sum":           {strconv.FormatInt(price, 10)},
		"receiver":      {s.payment.Wallet},
		"quickpay-form": {"donate"},
		"label":         {uEnc},
	}

	resURL := base.String() + "?" + params.Encode()

	log.WithFields(log.Fields{
		"chatID": chatID,
		"url":    resURL,
		"uEnc":   uEnc,
	}).Info("GenerateURL")
	return resURL
}

func (s *Service) GenerateProduct(countBought int64, price int64) string {
	return fmt.Sprintf(s.payment.Message, countBought, price)
}

func (s *Service) ProcessPayment(uEnc string, paidSum string) {
	log.WithFields(log.Fields{
		"uEnc": uEnc,
	}).Info("ProcessPayment")

	chatID, countBought, err := s.decryptProductAES(uEnc)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("decryptProductAES")
		return
	}

	paidSumDigits := strings.Split(paidSum, ".")
	paidSumInt, err := strconv.ParseInt(paidSumDigits[0], 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("ParseInt paidSum")
		return
	}

	m := make(map[int64]int64)

	m[s.payment.Products.ProductCount01] = s.payment.Products.ProductPrice01
	m[s.payment.Products.ProductCount02] = s.payment.Products.ProductPrice02
	m[s.payment.Products.ProductCount03] = s.payment.Products.ProductPrice03

	if paidSumInt != m[countBought] {
		log.WithFields(log.Fields{
			"paidSumInt": paidSumInt,
		}).Error("paidSumInt error")
		return
	}

	err = s.UpdateSubscription(chatID, countBought)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("UpdateBalance")
		return
	}
}
