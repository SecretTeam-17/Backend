package swagger

import (
	"log/slog"
	"net/http"
	"petsittersGameServer/swagger"

	"github.com/go-chi/chi/v5/middleware"
)

// New - возвращает новый обработчик для страницы swagger.
func New(alog slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.swagger.New"

		log := &alog
		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to receive swagger page")
		log = nil

		http.FileServer(swagger.Serve()).ServeHTTP(w, r)
	}
}
