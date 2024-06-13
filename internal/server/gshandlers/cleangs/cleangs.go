package cleangs

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/logger"
	rp "petsittersGameServer/internal/server/gshandlers/response"
	"petsittersGameServer/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type SessionCleaner interface {
	CleanSession(ctx context.Context, id string) (*storage.GameSession, error)
}

// New - возвращает новый хэндлер для очистки игровой сессии по id.
func New(alog slog.Logger, st SessionCleaner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.cleangs.New"

		log := &alog
		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to clean a game session")

		// Получаем параметр из запроса.
		id := chi.URLParam(r, "id")
		ctx := r.Context()

		// Получаем очищенную игровую сессию из БД.
		gs, err := st.CleanSession(ctx, id)
		if errors.Is(err, storage.ErrInput) {
			log.Error("invalid game session id", logger.Err(err))
			render.Status(r, 400)
			render.PlainText(w, r, "Error, failed to receive a game session: incorrect id")
			return
		}
		if errors.Is(err, storage.ErrSessionNotFound) {
			log.Error("game session not found", slog.String("id", id))
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to clean a game session: id not found")
			return
		}
		if err != nil {
			log.Error("failed to clean a game session", slog.String("id", id), logger.Err(err))
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to clean a game session: unknown error")
			return
		}
		log.Info("game session was cleaned successfully", slog.String("id", id))

		// Записываем сессию в структуру Response.
		var resp rp.Response
		resp.GameSession = *gs
		render.Status(r, 200)
		render.JSON(w, r, resp)
		log = nil
	}
}
