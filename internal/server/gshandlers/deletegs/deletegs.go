package deletegs

import (
	"errors"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/logger"
	"petsittersGameServer/internal/storage"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type SessionDeleter interface {
	DeleteSessionById(id int) error
}

// New - возвращает новый хэндлер для удаления игровой сессии по id.
func New(log *slog.Logger, st SessionDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.deletegs.New"

		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to delete a game session by id")

		// Получаем параметр из запроса и приводим его к типу int
		param := chi.URLParam(r, "id")
		id, err := strconv.Atoi(param)
		if err != nil {
			log.Error("invalid game session id", logger.Err(err))
			w.WriteHeader(400)
			render.PlainText(w, r, "Error, failed to delete a game session: incorrect id")
			return
		}

		// Удаляем игровую сессию из БД
		err = st.DeleteSessionById(id)
		if errors.Is(err, storage.ErrSessionNotFound) {
			log.Error("game session not found", slog.Int("session_id", id))
			w.WriteHeader(404)
			render.PlainText(w, r, "Error, failed to delete a game session: id not found")
			return
		}
		if err != nil {
			log.Error("failed to delete a game session", slog.Int("session_id", id))
			w.WriteHeader(404)
			render.PlainText(w, r, "Error, failed to delete a game session: unknown error")
			return
		}
		log.Info("game session was deleted successfully", slog.Int("session_id", id))

		// Возвращаем статус 204 и пустое тело
		render.Status(r, 204)
		render.NoContent(w, r)
	}
}
