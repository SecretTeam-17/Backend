package creategs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/logger"
	rp "petsittersGameServer/internal/server/gshandlers/response"
	"petsittersGameServer/internal/storage"
	"petsittersGameServer/internal/tools/api"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// Request - структура запроса для создания игровой сессии.
type Request struct {
	Name      string          `json:"username" validate:"required,max=100"`
	Email     string          `json:"email" validate:"required,email,max=100"`
	Stats     json.RawMessage `json:"stats"`
	Modules   json.RawMessage `json:"modules"`
	Minigames json.RawMessage `json:"minigames"`
}

type SessionCreator interface {
	CreateSession(ctx context.Context, name, email string, stats, modules, minigames []byte) (*storage.GameSession, error)
}

// New - возвращает новый хэндлер для создания игровой сессии.
func New(alog slog.Logger, st SessionCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.creategs.New"

		log := &alog
		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to create a game session")

		// Декодируем тело запроса в структуру Request и проверяем на ошибки.
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.Status(r, 400)
			render.PlainText(w, r, "Error, failed to create new gameSession: empty request")
			return
		}
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			render.Status(r, 400)
			render.PlainText(w, r, "Error, failed to create new gameSession: failed to decode request")
			return
		}
		log.Info("request body decoded")

		// // Преобразуем поля stats, modules и minigames в слайсы байт.
		// stats, err := req.Stats.MarshalJSON()
		// if err != nil {
		// 	log.Error("failed to marshal JSON stats", logger.Err(err))
		// 	render.Status(r, 400)
		// 	render.PlainText(w, r, "Error, failed to create new gameSession: failed to marshal JSON stats")
		// 	return
		// }
		// modules, err := req.Modules.MarshalJSON()
		// if err != nil {
		// 	log.Error("failed to marshal JSON modules", logger.Err(err))
		// 	render.Status(r, 400)
		// 	render.PlainText(w, r, "Error, failed to create new gameSession: failed to marshal JSON modules")
		// 	return
		// }
		// minigames, err := req.Minigames.MarshalJSON()
		// if err != nil {
		// 	log.Error("failed to marshal JSON minigames", logger.Err(err))
		// 	render.Status(r, 400)
		// 	render.PlainText(w, r, "Error, failed to create new gameSession: failed to marshal JSON minigames")
		// 	return
		// }

		// Валидация полей json из запроса.
		valid := validator.New()
		err = valid.Struct(req)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", logger.Err(err))
			render.Status(r, 422)
			str := fmt.Sprintf("Error, failed to create new gameSession: %s", api.ValidationError(validateErr))
			render.PlainText(w, r, str)
			return
		}
		req.Email = strings.ToLower(req.Email)

		ctx := r.Context()

		// Создаем нового юзера и игровую сессию по данным из запроса.
		gs, err := st.CreateSession(ctx, req.Name, req.Email, req.Stats, req.Modules, req.Minigames)
		// Если игрок с данным email уже существует, то возвращаем его игровую сессию.
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("game session already exists; returning session data", slog.String("email", req.Email))
			render.Status(r, 200)
			render.JSON(w, r, gs)
			return
		}
		if errors.Is(err, storage.ErrInput) {
			log.Error("incorrect input user data", slog.String("name", req.Name), slog.String("email", req.Email))
			render.Status(r, 422)
			render.PlainText(w, r, "Error, failed to create new gameSession: incorrect input user data")
			return
		}
		if err != nil {
			log.Error("failed to create gameSession", logger.Err(err))
			render.Status(r, 422)
			render.PlainText(w, r, "Error, failed to create new gameSession: unknown error")
			return
		}
		log.Info("new gameSession created", slog.String("id", gs.Id.Hex()))

		// Записываем данные сессии в структуру Response
		var resp rp.Response
		resp.GameSession = *gs
		render.Status(r, 201)
		render.JSON(w, r, resp)
		log = nil
	}
}
