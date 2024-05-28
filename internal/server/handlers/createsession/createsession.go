package createsession

import (
	"errors"
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

// Request - структура запроса для создания игровой сессии.
type Request struct {
	Name  string `json:"username"`
	Email string `json:"email" validate:"required,email"`
}

// Response - структура ответа со статусом, ошибкой и созданной сессией.
type Response struct {
	api.RespStatus
	GameSession storage.GameSession `json:"gameSession"`
}

// Возможно, интерфейсы хранилища лучше перенести в пакет storage
type SessionCreator interface {
	CreateSession(name, email string) (*storage.GameSession, error)
}

// New - возвращает новый хэндлер для создания игровой сессии.
func New(log *slog.Logger, sc SessionCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.createsession.New"

		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Декодируем тело запроса в структуру Request и проверяем на ошибки
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, api.Error("failed to create new gameSession, empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			render.JSON(w, r, api.Error("failed to create new gameSession, failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		valid := validator.New()
		err = valid.Struct(req)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", logger.Err(err))

			// render.JSON(w, r, api.Error("invalid request"))
			render.JSON(w, r, api.ValidationError(validateErr))
			return
		}

		// Создаем нового юзера и игровую сессию по данным из запроса
		gs, err := sc.CreateSession(req.Name, req.Email)
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("user already exists", slog.String("email", req.Email))
			render.JSON(w, r, api.Error("failed to create new gameSession, user already exists"))
			return
		}
		if errors.Is(err, storage.ErrInput) {
			log.Error("incorrect input user data", slog.String("name", req.Name), slog.String("email", req.Email))
			render.JSON(w, r, api.Error("failed to create new gameSession, incorrect input user data"))
			return
		}
		if err != nil {
			log.Error("failed to create gameSession", logger.Err(err))
			render.JSON(w, r, api.Error("failed to create new gameSession"))
			return
		}
		log.Info("new gameSession created", slog.Int("id", gs.SessionID))

		// Записывает данные юзера и сессии в структуру Response
		var resp Response
		resp.GameSession = *gs
		resp.RespStatus = api.OK()
		render.JSON(w, r, resp)
	}
}
