package service

import (
	log "github.com/sirupsen/logrus"
)

func (s *Service) GetSales() int64 {
	log.Info("GetSales")

	return s.payment.Products.SalesCount
}

func (s *Service) SetSales(count int64)  {
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

	s.payment.Products.SalesCount = count

	log.WithFields(log.Fields{
		"sales now": count,
	}).Info("SetSales")
}


