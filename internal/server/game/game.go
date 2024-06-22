package game

import (
	"net/http"
	"petsittersGameServer/game"
)

// New - возвращает новый обработчик для страницы игры.
func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(game.Serve()).ServeHTTP(w, r)
	}
}
