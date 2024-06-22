package swagger

import (
	"net/http"
	"petsittersGameServer/swagger"
)

// New - возвращает новый обработчик для страницы swagger.
func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(swagger.Serve()).ServeHTTP(w, r)
	}
}
