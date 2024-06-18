package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"petsittersGameServer/internal/config"
	"petsittersGameServer/internal/logger"
	"petsittersGameServer/internal/server/gshandlers/cleangs"
	"petsittersGameServer/internal/server/gshandlers/creategs"
	"petsittersGameServer/internal/server/gshandlers/deletegs"
	"petsittersGameServer/internal/server/gshandlers/getallgs"
	"petsittersGameServer/internal/server/gshandlers/getgsemail"
	"petsittersGameServer/internal/server/gshandlers/getgsid"
	"petsittersGameServer/internal/server/gshandlers/truncate"
	"petsittersGameServer/internal/server/gshandlers/updategs"
	"petsittersGameServer/internal/server/index/indexpage"
	"petsittersGameServer/internal/server/middleware/cors"
	"petsittersGameServer/internal/server/middleware/redirect"
	"petsittersGameServer/internal/storage/mongodb"
	"petsittersGameServer/internal/tools/stopsignal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	// Загружаем данные из конфиг файла
	cfg := config.MustLoad()

	// Инициализируем и настраиваем основной логгер
	log := logger.SetupLogger(cfg.Env)
	log.Debug("logger initialized")

	// Инициализируем пул подключений к базе данных
	//storage, err := sqlite.New(cfg.StoragePath, cfg.ModulesPath)
	storage, err := mongodb.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", logger.Err(err))
		os.Exit(1)
	}
	log.Debug("storage initialized")

	// Получаем главный роутер
	router := chi.NewRouter()

	// Указываем, какие middleware использовать
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Cors())

	router.Get("/", indexpage.New(*log))

	// Объявляем REST API хэндлеры для работы со структурой GameSession
	router.Get("/api/session/id/{id}", getgsid.New(*log, storage))
	router.Get("/api/session/email/{email}", getgsemail.New(*log, storage))
	router.Get("/api/session/all", getallgs.New(*log, storage))
	router.Get("/api/session/new/{id}", cleangs.New(*log, storage))

	router.Post("/api/session", creategs.New(*log, storage))

	router.Put("/api/session", updategs.New(*log, storage))

	router.Delete("/api/session/id/{id}", deletegs.New(*log, storage))
	router.Delete("/api/session/verydangerousbutton", truncate.New(*log, storage)) // Для тестирования

	// Конфигурируем сервер из данных конфиг файла
	srv := &http.Server{
		Addr:         cfg.AddressTLS,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	log.Info("starting server", slog.String("env", cfg.Env), slog.String("address", cfg.Address))

	// Запускаем сервер TLS в отдельной горутине
	go func() {
		if err := srv.ListenAndServeTLS(cfg.CertPath, cfg.CertKeyPath); err != nil {
			log.Error("failed to start server with TLS")
		}
	}()
	log.Info("server with TLS started")

	// Запускаем сервер с редиректом в отдельной горутине
	go func() {
		if err := http.ListenAndServe(cfg.Address, http.HandlerFunc(redirect.ToTLS)); err != nil {
			log.Error("failed to start redirect server")
		}
	}()
	log.Info("redirect server started")

	// Блокируем выполнение основной горутины
	log.Info("server awaits INT signal to stop")
	stopsignal.Stop()

	log.Info("stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Останавливаем сервер
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", logger.Err(err))
		os.Exit(1)
	}
	// Закрываем базу данных
	storage.Close()

	log.Info("server stopped")
}
