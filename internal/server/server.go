package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"petsittersGameServer/internal/config"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	srv  *http.Server
	cert string
	key  string
}

// New - создает новый HTTPS сервер.
func New(cfg *config.Config, router chi.Router) *Server {
	server := &Server{
		srv: &http.Server{
			Addr:         cfg.AddressTLS,
			Handler:      router,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		cert: cfg.CertPath,
		key:  cfg.CertKeyPath,
	}
	return server
}

// Start - запускает HTTPS сервер в отдельной горутине.
func (s *Server) Start() {
	go func() {
		if err := s.srv.ListenAndServeTLS(s.cert, s.key); err != nil {
			// Игнорируем ошибку при вызове Shutdown.
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Printf("failed to start server with TLS: %s", err.Error())
		}
	}()
}

// Shutdown - останавливает сервер используя graceful shutdown.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to stop server: %s", err.Error())
	}
}
