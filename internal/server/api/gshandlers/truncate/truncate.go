package truncate

import (
	"errors"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/logger"
	"petsittersGameServer/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// type DataTruncater interface {
// 	TruncateData(ctx context.Context) error
// }

// New - возвращает новый хэндлер для удаления всех данных из таблиц игроков и игровых сессий.
func New(alog slog.Logger, st storage.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.truncate.New"

		log := &alog
		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to truncate tables")

		ctx := r.Context()

		// Выполняем очистку таблиц.
		err := st.TruncateData(ctx)
		if errors.Is(err, storage.ErrSessionsEmpty) {
			log.Error("game sessions not found", logger.Err(err))
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to delete a game session: no sessions found")
			return
		}
		if err != nil {
			log.Error("failed to truncate tables", logger.Err(err))
			render.Status(r, 500)
			render.PlainText(w, r, "Error, failed to truncate collection: unknown error")
			return
		}
		log.Info("collection was truncated successfully")

		// Возвращаем статус 204 и пустое тело.
		render.Status(r, 204)
		render.NoContent(w, r)
		log = nil
	}
}
