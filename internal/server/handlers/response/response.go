package response

import (
	"petsittersGameServer/internal/storage"
)

// Response - структура ответа со статусом, ошибкой и созданной сессией.
type Response struct {
	storage.GameSession
}
