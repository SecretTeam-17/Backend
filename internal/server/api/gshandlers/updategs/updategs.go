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

// type SessionUpdater interface {
// 	UpdateSession(ctx context.Context, gs storage.GameSession) error
// }

// New - возвращает новый хэндлер для обновления игровой сессии.
func New(alog slog.Logger, st storage.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.updategs.New"

		log := &alog
		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to update a game session")

		// Декодируем тело запроса в структуру GameSession и проверяем на ошибки.
		var req storage.GameSession
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.Status(r, 400)
			render.PlainText(w, r, "Error, failed to update a game session: empty request")
			return
		}
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			render.Status(r, 400)
			render.PlainText(w, r, "Error, failed to update a game session: failed to decode request; incorrect input")
			return
		}
		log.Info("request body decoded")

		// Валидация полей json из запроса.
		valid := validator.New()
		err = valid.Struct(req)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", logger.Err(err))
			render.Status(r, 422)
			str := fmt.Sprintf("Error, failed to update a game session: %s", api.ValidationError(validateErr))
			render.PlainText(w, r, str)
			return
		}

		ctx := r.Context()

		// Обновляем игровую сессию в БД.
		err = st.UpdateSession(ctx, req)
		if errors.Is(err, storage.ErrSessionNotFound) {
			log.Error("game session not found", slog.String("id", req.Id.Hex()))
			render.Status(r, 422)
			render.PlainText(w, r, "Error, to update a game session: id not found")
			return
		}
		if err != nil {
			log.Error("failed to update a game session", logger.Err(err))
			render.Status(r, 422)
			render.PlainText(w, r, "Error, failed to update a game session: unknown error")
			return
		}
		log.Info("game session updated", slog.String("id", req.Id.Hex()))

		// Возвращаем статус 204 и пустое тело.
		render.Status(r, 204)
		render.NoContent(w, r)
		log = nil
	}
}
