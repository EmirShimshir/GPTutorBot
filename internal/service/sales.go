package service

import (
	log "github.com/sirupsen/logrus"
)

func (s *Service) GetSales() (int64, error) {
	log.Info("GetSales")

	count, err := s.repo.Sales.Get()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Service) SetSales(count int64) error {
	log.Info("SetSales")

	if count > 0 {
		s.payment.Products.ProductPrice01 = s.payment.Products.ProductSalesPrice01
		s.payment.Products.ProductPrice02 = s.payment.Products.ProductSalesPrice02
		s.payment.Products.ProductPrice03 = s.payment.Products.ProductSalesPrice03
	} else {
		s.payment.Products.ProductPrice01 = s.payment.Products.ProductBasePrice01
		s.payment.Products.ProductPrice02 = s.payment.Products.ProductBasePrice02
		s.payment.Products.ProductPrice03 = s.payment.Products.ProductBasePrice03
	}

	err := s.repo.Sales.Save(count)
	if err != nil {
		return err
	}

	return nil
}


