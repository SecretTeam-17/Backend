package getgsemail

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"petsittersGameServer/internal/logger"
	rp "petsittersGameServer/internal/server/gshandlers/response"
	"petsittersGameServer/internal/storage"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type SessionByEmail interface {
	GetSessionByEmail(ctx context.Context, email string) (*storage.GameSession, error)
}

// New - возвращает новый хэндлер для получения игровой сессии по email.
func New(ctx context.Context, alog slog.Logger, st SessionByEmail) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.getgsemail.New"

		log := &alog
		log = log.With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("new request to receive a game session by email")

		// Получаем параметр из запроса и валидируем его к формату электронной почты
		email := chi.URLParam(r, "email")
		email, err := url.QueryUnescape(email)
		if err != nil {
			log.Error("failed to decode email in URL")
			render.Status(r, 400)
			render.PlainText(w, r, "Error, failed to decode email in URL: incorrect email")
			return
		}
		email = strings.ToLower(email)
		valid := validator.New()
		err = valid.Var(email, "required,email")
		if err != nil {
			log.Error("invalid user email", logger.Err(err))
			render.Status(r, 422)
			render.PlainText(w, r, "Error, failed to receive a game session: incorrect email")
			return
		}

		// Получаем игровую сессию из БД по email ее юзера
		gs, err := st.GetSessionByEmail(ctx, email)
		if errors.Is(err, storage.ErrSessionNotFound) {
			log.Error("game session not found", slog.String("user_email", email))
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to receive a game session: user email not found")
			return
		}
		if err != nil {
			log.Error("failed to receive a game session", slog.String("user_email", email), logger.Err(err))
			render.Status(r, 404)
			render.PlainText(w, r, "Error, failed to receive a game session: unknown error")
			return
		}
		log.Info("game session was found successfully", slog.String("user_email", email))

		// Записываем сессию в структуру Response
		var resp rp.Response
		resp.GameSession = *gs
		render.Status(r, 200)
		render.JSON(w, r, resp)
		log = nil
	}
}
