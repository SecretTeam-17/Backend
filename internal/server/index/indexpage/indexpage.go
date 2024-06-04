package indexpage

import (
	"log/slog"
	"net/http"
	"os"
	"petsittersGameServer/internal/logger"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// New - возвращает новый хэндлер для index page.
func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.indexpage.New"

		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to receive index page")

		buf, err := os.ReadFile("./internal/templates/index.html")
		if err != nil {
			log.Error("can't read index.html", logger.Err(err))
			w.WriteHeader(500)
			render.PlainText(w, r, "Error, failed to read index.html")
			return
		}
		render.Status(r, 200)
		render.HTML(w, r, string(buf))
	}
}
