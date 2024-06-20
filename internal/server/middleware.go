package server

import (
	"petsittersGameServer/internal/server/middleware/cors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Middleware - инициализирует все необходимые обработчики middleware.
func Middleware(router chi.Router) {
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Cors())
}
