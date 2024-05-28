// Пакет для использования в тестах.
package slogdiscard

import (
	"context"
	"log/slog"
)

// NewDiscardLogger - создает новый логгер, который можно передать в тесты для отмены записи этих тестов в логи.
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler - хэндлер отмены записи логов
type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Описываем 4 метода, чтобы удовлетворять интерфейсу Handler
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
