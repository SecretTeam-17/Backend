package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"petsittersGameServer/internal/config"
	"petsittersGameServer/internal/logger"
	"petsittersGameServer/internal/server/handlers/createsession"
	"petsittersGameServer/internal/storage/sqlite"
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
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", logger.Err(err))
		os.Exit(1)
	}
	log.Debug("storage initialized")

	// Получаем главный роутер
	router := chi.NewRouter()

	// Указываем, какие middleware использовать
	router.Use(middleware.RequestID)
	// router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// router.Use(middleware.URLFormat)

	// Объявляем REST хэндлеры
	router.Post("/new", createsession.New(log, storage))

	// Конфигурируем сервер из данных конфиг файла
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	log.Info("starting server", slog.String("env", cfg.Env), slog.String("address", cfg.Address))

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()
	log.Info("server started")

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
