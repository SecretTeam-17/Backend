package server

import (
	"log/slog"
	"petsittersGameServer/internal/server/game"
	"petsittersGameServer/internal/tools/static"

	"github.com/go-chi/chi/v5"
)

// Game - инициализирует обработчик для страницы игры.
func Game(router chi.Router, log *slog.Logger) {
	static.FileServer(router, "/", game.New(*log))
}
