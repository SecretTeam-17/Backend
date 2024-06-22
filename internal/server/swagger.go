package server

import (
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/server/swagger"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Swagger - инициализирует обработчик для страницы swagger.
func Swagger(router chi.Router, log *slog.Logger) {
	const operation = "server.Swagger"

	path := "/swagger"
	router.Get(path, http.RedirectHandler("/api/swagger/", http.StatusMovedPermanently).ServeHTTP)

	path += "/*"
	router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/swagger/" {
			log.Info("new request to receive swagger page", slog.String("op", operation))
		}
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, swagger.New())
		fs.ServeHTTP(w, r)
	})
}
