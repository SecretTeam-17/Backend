package truncate

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type DataTruncater interface {
	TruncateTables() error
}

// New - возвращает новый хэндлер для удаления всех данных из таблиц игроков и игровых сессий.
func New(log *slog.Logger, st DataTruncater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.truncate.New"

		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to truncate tables")

		// Выполняем очистку таблиц
		err := st.TruncateTables()
		if err != nil {
			fmt.Println(err)
			log.Error("failed to truncate tables")
			w.WriteHeader(500)
			render.PlainText(w, r, "Error, failed to truncate tables: unknown error")
			return
		}
		log.Info("tables was truncated successfully")

		// Возвращаем статус 204 и пустое тело
		w.WriteHeader(204)
		render.NoContent(w, r)
	}
}