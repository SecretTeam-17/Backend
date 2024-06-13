package deletegs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/logger"
	"petsittersGameServer/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type SessionDeleter interface {
	DeleteSessionById(ctx context.Context, id string) error
}

// New - возвращает новый хэндлер для удаления игровой сессии по id.
func New(alog slog.Logger, st SessionDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.deletegs.New"

		log := &alog
		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to delete a game session by id")

		// Получаем параметр id из запроса.
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		// Удаляем игровую сессию из БД.
		err := st.DeleteSessionById(ctx, id)
		if errors.Is(err, storage.ErrInput) {
			log.Error("invalid game session id", logger.Err(err))
			render.Status(r, 400)
			render.PlainText(w, r, "Error, failed to receive a game session: incorrect id")
			return
		}
		if errors.Is(err, storage.ErrSessionNotFound) {
			log.Error("game session not found", slog.String("id", id))
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to delete a game session: id not found")
			return
		}
		if err != nil {
			log.Error("failed to delete a game session", slog.String("id", id), logger.Err(err))
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to delete a game session: unknown error")
			return
		}
		log.Info("game session was deleted successfully", slog.String("id", id))

		// Возвращаем статус 204 и пустое тело.
		render.Status(r, 204)
		render.NoContent(w, r)
		log = nil
	}
}
