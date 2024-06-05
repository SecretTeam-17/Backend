package getgsid

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/logger"
	rp "petsittersGameServer/internal/server/gshandlers/response"
	"petsittersGameServer/internal/storage"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type SessionById interface {
	GetSessionById(ctx context.Context, id int) (*storage.GameSession, error)
}

// New - возвращает новый хэндлер для получения игровой сессии по id.
func New(ctx context.Context, log *slog.Logger, st SessionById) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.getgsid.New"

		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to receive a game session by id")

		// Получаем параметр из запроса и приводим его к типу int
		param := chi.URLParam(r, "id")
		id, err := strconv.Atoi(param)
		if err != nil {
			log.Error("invalid game session id", logger.Err(err))
			render.Status(r, 400)
			render.PlainText(w, r, "Error, failed to receive a game session: incorrect id")
			return
		}

		// Получаем игровую сессию из БД по ее id
		gs, err := st.GetSessionById(ctx, id)
		if errors.Is(err, storage.ErrSessionNotFound) {
			log.Error("game session not found", slog.Int("session_id", id))
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to receive a game session: id not found")
			return
		}
		if err != nil {
			log.Error("failed to receive a game session", slog.Int("session_id", id), logger.Err(err))
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to receive a game session: unknown error")
			return
		}
		log.Info("game session was found successfully", slog.Int("session_id", id))

		// Записываем сессию в структуру Response
		var resp rp.Response
		resp.GameSession = *gs
		render.Status(r, 200)
		render.JSON(w, r, resp)
	}
}
