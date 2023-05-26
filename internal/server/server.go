package server

import (
	"fmt"
	"github.com/EmirShimshir/tasker-bot/internal/config"
	"github.com/EmirShimshir/tasker-bot/internal/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	server   *http.Server
	services *service.Service
	port     uint
}

func NewPayment(services *service.Service, cfg config.Server) *Server {
	return &Server{
		services: services,
		port:     cfg.Port,
	}
}

func (s *Server) Start() error {

	s.server = &http.Server{
		Handler: s,
		Addr:    fmt.Sprintf(":%d", s.port),
	}

	return s.server.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.URL.Path == "/payment/umoney" {
		uEnc := r.FormValue("label")
		paidSum := r.FormValue("withdraw_amount")

		log.WithFields(log.Fields{
			"uEnc": uEnc,
		}).Info("new server payment handle")

		s.services.ProcessPayment(uEnc, paidSum)

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	return
}
