package server

import (
	"log/slog"
	"petsittersGameServer/internal/server/swagger"
	"petsittersGameServer/internal/tools/static"

	"github.com/go-chi/chi/v5"
)

// Swagger - инициализирует обработчик для страницы swagger.
func Swagger(router chi.Router, log *slog.Logger) {
	static.FileServer(router, "/api/swagger", swagger.New(*log))
}
