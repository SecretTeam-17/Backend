package server

import (
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/server/game"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Game - инициализирует обработчик для страницы игры.
func Game(router chi.Router, log *slog.Logger) {
	const operation = "server.Game"

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			log.Info("new request to receive game page", slog.String("op", operation))
		}
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, game.New())
		fs.ServeHTTP(w, r)
	})
}
