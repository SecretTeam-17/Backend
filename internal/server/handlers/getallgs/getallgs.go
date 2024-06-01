package getallgs

import (
	"errors"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type AllSessions interface {
	GetSessions() ([]storage.GameSession, error)
}

// New - возвращает новый хэндлер для получения всех игровых сессий.
func New(log *slog.Logger, st AllSessions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.getallgs.New"

		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to receive all game sessions")

		// Получаем слайс всех игровых сессий из БД
		resp, err := st.GetSessions()
		if errors.Is(err, storage.ErrSessionsEmpty) {
			log.Error("game sessions not found")
			w.WriteHeader(404)
			render.PlainText(w, r, "Error, failed to receive all game session: table is empty")
			return
		}
		if err != nil {
			log.Error("failed to receive all game sessions")
			w.WriteHeader(404)
			render.PlainText(w, r, "Error, failed to receive all game sessions: unknown error")
			return
		}

		// Записываем слайс игровых сессий из БД в респонс
		w.WriteHeader(200)
		render.JSON(w, r, resp)
	}
}
