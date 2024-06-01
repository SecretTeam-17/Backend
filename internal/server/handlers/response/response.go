package response

import (
	"petsittersGameServer/internal/storage"
)

// Response - структура ответа с игровой сессией.
type Response struct {
	storage.GameSession
}
