package response

import (
	"petsittersGameServer/internal/storage"
	"petsittersGameServer/internal/storage/sqlite"
)

// ResponseSQL - структура ответа с игровой сессией.
type ResponseSQL struct {
	sqlite.GameSession
}

// Response - структура ответа с игровой сессией.
type Response struct {
	storage.GameSession
}
