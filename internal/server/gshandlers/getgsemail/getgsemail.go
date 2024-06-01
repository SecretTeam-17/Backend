package getgsemail

import (
	"errors"
	"log/slog"
	"net/http"
	"petsittersGameServer/internal/logger"
	rp "petsittersGameServer/internal/server/gshandlers/response"
	"petsittersGameServer/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type SessionByEmail interface {
	GetSessionByEmail(email string) (*storage.GameSession, error)
}

// New - возвращает новый хэндлер для получения игровой сессии по email.
func New(log *slog.Logger, st SessionByEmail) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.getgsemail.New"

		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to receive a game session by email")

		// Получаем параметр из запроса и валидируем его к формату электронной почты
		email := chi.URLParam(r, "email")
		valid := validator.New()
		err := valid.Var(email, "required,email")
		if err != nil {
			log.Error("invalid user email", logger.Err(err))
			w.WriteHeader(422)
			render.PlainText(w, r, "Error, failed to receive a game session: incorrect email")
			return
		}

		// Получаем игровую сессию из БД по email ее юзера
		gs, err := st.GetSessionByEmail(email)
		if errors.Is(err, storage.ErrSessionNotFound) {
			log.Error("game session not found", slog.String("user_email", email))
			w.WriteHeader(404)
			render.PlainText(w, r, "Error, failed to receive a game session: user email not found")
			return
		}
		if err != nil {
			log.Error("failed to receive a game session", slog.String("user_email", email))
			w.WriteHeader(404)
			render.PlainText(w, r, "Error, failed to receive a game session: unknown error")
			return
		}
		log.Info("game session was found successfully", slog.String("user_email", email))

		// Записываем сессию в структуру Response
		var resp rp.Response
		resp.GameSession = *gs
		w.WriteHeader(200)
		render.JSON(w, r, resp)
	}
}
