package updategs

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/logger"
	"petsittersGameServer/internal/storage"
	"petsittersGameServer/internal/tools/api"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type SessionUpdater interface {
	UpdateSession(gs storage.GameSession) error
}

// New - возвращает новый хэндлер для обновления игровой сессии.
func New(log *slog.Logger, st SessionUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.updategs.New"

		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to update a game session")

		// Декодируем тело запроса в структуру GameSession и проверяем на ошибки
		var req storage.GameSession
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			w.WriteHeader(400)
			render.PlainText(w, r, "Error, failed to update a game session: empty request")
			return
		}
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			w.WriteHeader(400)
			render.PlainText(w, r, "Error, failed to update a game session: failed to decode request")
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		// Валидация полей json из запроса
		valid := validator.New()
		err = valid.Struct(req)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", logger.Err(err))
			w.WriteHeader(422)
			str := fmt.Sprintf("Error, failed to update a game session: %s", api.ValidationError(validateErr))
			render.PlainText(w, r, str)
			return
		}

		// Обновляем игровую сессию в БД
		err = st.UpdateSession(req)
		if errors.Is(err, storage.ErrModuleNotFound) {
			log.Error("incorrect module", slog.Int("module", req.CurrentModule))
			w.WriteHeader(422)
			render.PlainText(w, r, "Error, to update a game session: module not found in database")
			return
		}
		if errors.Is(err, storage.ErrSessionNotFound) {
			log.Error("game session not found", slog.Int("session_id", req.SessionID))
			w.WriteHeader(404)
			render.PlainText(w, r, "Error, to update a game session: id not found")
			return
		}
		if err != nil {
			log.Error("failed to update a game session", logger.Err(err))
			w.WriteHeader(422)
			render.PlainText(w, r, "Error, failed to update a game session: unknown error")
			return
		}
		log.Info("game session updated", slog.Int("id", req.SessionID))

		// Возвращаем статус 204 и пустое тело
		render.Status(r, 204)
		render.NoContent(w, r)
	}
}
