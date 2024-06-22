package game

import (
	"log/slog"
	"net/http"
	"petsittersGameServer/game"

	"github.com/go-chi/chi/v5/middleware"
)

// New - возвращает новый обработчик для страницы игры.
func New(alog slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.game.New"

		log := &alog
		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to receive the game")
		log = nil

		http.FileServer(game.Serve()).ServeHTTP(w, r)
	}
}
