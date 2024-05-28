package logger

import (
	"log/slog"
	"os"
)

// Константы окружений
const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// SetupLogger - задает вывод логов в виде текста либо JSON в зависимости от окружения.
func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // По дефолту выбирается prod режим логгера
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	// TODO: добавить логирование в файл
	return log
}

// Err - представляет ошибку как атрибут слоггера.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
