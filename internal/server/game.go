package server

import (
	"log/slog"
	"petsittersGameServer/internal/server/game"

	"github.com/go-chi/chi/v5"
)

// Game - инициализирует обработчик запроса игры.
func Game(router chi.Router, log *slog.Logger) {
	router.Handle("/", game.New(*log))
}
