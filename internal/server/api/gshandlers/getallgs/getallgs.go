package getallgs

import (
	"errors"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/logger"
	"petsittersGameServer/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// type AllSessions interface {
// 	GetSessions(ctx context.Context) ([]storage.GameSession, error)
// }

// New - возвращает новый хэндлер для получения всех игровых сессий.
func New(alog slog.Logger, st storage.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.getallgs.New"

		log := &alog
		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to receive all game sessions")

		ctx := r.Context()

		// Получаем слайс всех игровых сессий из БД.
		resp, err := st.GetSessions(ctx)
		if errors.Is(err, storage.ErrSessionsEmpty) {
			log.Error("game sessions not found")
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to receive all game sessions: no sessions found")
			return
		}
		if err != nil {
			log.Error("failed to receive all game sessions", logger.Err(err))
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to receive all game sessions: unknown error")
			return
		}
		log.Info("all game sessions was returned successfully")

		// Записываем слайс игровых сессий из БД в респонс.
		render.Status(r, 200)
		render.JSON(w, r, resp)
		log = nil
	}
}
