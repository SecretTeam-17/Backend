package server

import (
	"log/slog"
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
)

// API - инициализирует все необходимые обработчики API.
func API(router chi.Router, log *slog.Logger, storage storage.Interface) {
	router.Get("/api/session/id/{id}", getgsid.New(*log, storage))
	router.Get("/api/session/email/{email}", getgsemail.New(*log, storage))
	router.Get("/api/session/all", getallgs.New(*log, storage))
	router.Get("/api/session/new/{id}", cleangs.New(*log, storage))

	router.Post("/api/session", creategs.New(*log, storage))

	router.Put("/api/session", updategs.New(*log, storage))

	router.Delete("/api/session/id/{id}", deletegs.New(*log, storage))
	router.Delete("/api/session/verydangerousbutton", truncate.New(*log, storage)) // Для тестирования
}
