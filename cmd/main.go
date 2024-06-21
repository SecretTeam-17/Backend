package main

import (
	"log/slog"
	"petsittersGameServer/internal/config"
	"petsittersGameServer/internal/logger"
	"petsittersGameServer/internal/server"
	"petsittersGameServer/internal/server/redirect"
	"petsittersGameServer/internal/storage/mongodb"
	"petsittersGameServer/internal/tools/stopsignal"

	"github.com/go-chi/chi/v5"
)

func main() {

	// Загружаем данные из конфиг файла.
	cfg := config.MustLoad()

	// Инициализируем и настраиваем основной логгер.
	log := logger.SetupLogger(cfg.Env)
	log.Debug("logger initialized")

	// Инициализируем пул подключений к базе данных.
	storage := mongodb.New(cfg)
	log.Debug("storage initialized")
	defer storage.Close()

	// Получаем главный роутер.
	router := chi.NewRouter()

	// Указываем обработчики.
	server.Middleware(router)
	server.Game(router, log)
	server.API(router, log, storage)
	server.Swagger(router, log)

	// Создаем новый сервер и запускаем в отдельной горутине.
	srv := server.New(cfg, router)
	srv.Start()
	log.Info("server started", slog.String("env", cfg.Env), slog.String("address", cfg.Address))

	// Запускаем редирект на https в отдельной горутине.
	redirect.ToTLS(cfg)
	log.Info("redirect started")

	// Блокируем выполнение основной горутины.
	log.Info("server awaits INT signal to stop")
	stopsignal.Stop()

	// Остановливаем сервер после сигнала прерывания.
	srv.Shutdown()

	log.Info("server stopped")
}
