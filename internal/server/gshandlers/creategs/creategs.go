package creategs

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/logger"
	rp "petsittersGameServer/internal/server/gshandlers/response"
	"petsittersGameServer/internal/storage"
	"petsittersGameServer/internal/tools/api"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// Request - структура запроса для создания игровой сессии.
type Request struct {
	Name  string `json:"username" validate:"required,max=100"`
	Email string `json:"email" validate:"required,email,max=100"`
}

// Возможно, интерфейсы хранилища лучше перенести в пакет storage
type SessionCreator interface {
	CreateSession(name, email string) (*storage.GameSession, error)
}

// New - возвращает новый хэндлер для создания игровой сессии.
func New(log *slog.Logger, st SessionCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.creategs.New"

		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to create a game session")

		// Декодируем тело запроса в структуру Request и проверяем на ошибки
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			w.WriteHeader(400)
			render.PlainText(w, r, "Error, failed to create new gameSession: empty request")
			return
		}
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			w.WriteHeader(400)
			render.PlainText(w, r, "Error, failed to create new gameSession: failed to decode request")
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
			str := fmt.Sprintf("Error, failed to create new gameSession: %s", api.ValidationError(validateErr))
			render.PlainText(w, r, str)
			return
		}

		// Создаем нового юзера и игровую сессию по данным из запроса
		gs, err := st.CreateSession(req.Name, req.Email)
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("user already exists", slog.String("email", req.Email))
			w.WriteHeader(422)
			render.PlainText(w, r, "Error, failed to create new gameSession: user already exists")
			return
		}
		if errors.Is(err, storage.ErrInput) {
			log.Error("incorrect input user data", slog.String("name", req.Name), slog.String("email", req.Email))
			w.WriteHeader(422)
			render.PlainText(w, r, "Error, failed to create new gameSession: incorrect input user data")
			return
		}
		if err != nil {
			log.Error("failed to create gameSession", logger.Err(err))
			w.WriteHeader(422)
			render.PlainText(w, r, "Error, failed to create new gameSession: unknown error")
			return
		}
		log.Info("new gameSession created", slog.Int("id", gs.SessionID))

		// TODO: cookie
		// Записываем данные сессии в структуру Response
		var resp rp.Response
		resp.GameSession = *gs
		w.WriteHeader(201)
		render.JSON(w, r, resp)
	}
}