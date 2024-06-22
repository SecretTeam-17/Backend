package server

import (
	"log/slog"
	"petsittersGameServer/internal/config"
	"petsittersGameServer/internal/server/api/gshandlers/cleangs"
	"petsittersGameServer/internal/server/api/gshandlers/creategs"
	"petsittersGameServer/internal/server/api/gshandlers/deletegs"
	"petsittersGameServer/internal/server/api/gshandlers/getallgs"
	"petsittersGameServer/internal/server/api/gshandlers/getgsemail"
	"petsittersGameServer/internal/server/api/gshandlers/getgsid"
	"petsittersGameServer/internal/server/api/gshandlers/truncate"
	"petsittersGameServer/internal/server/api/gshandlers/updategs"
	"petsittersGameServer/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// API - инициализирует все общедоступные обработчики API.
func API(router chi.Router, log *slog.Logger, storage storage.Interface) {
	router.Get("/session/id/{id}", getgsid.New(*log, storage))
	router.Get("/session/email/{email}", getgsemail.New(*log, storage))
	router.Get("/session/new/{id}", cleangs.New(*log, storage))

	router.Post("/session", creategs.New(*log, storage))

	router.Put("/session", updategs.New(*log, storage))
}

// API - инициализирует все обработчики API, требующие авторизации.
func SensitiveAPI(cfg *config.Config, router chi.Router, log *slog.Logger, storage storage.Interface) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth("sensitive", map[string]string{
			cfg.AuthUser: cfg.AuthPasswd,
		}))

		r.Get("/session/all", getallgs.New(*log, storage))
		r.Delete("/session/id/{id}", deletegs.New(*log, storage))
		r.Delete("/session/verydangerousbutton", truncate.New(*log, storage)) // Для тестирования
	})
}
