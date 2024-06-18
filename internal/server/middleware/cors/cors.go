package cors

import (
	"net/http"

	"github.com/go-chi/cors"
)

func Cors() func(next http.Handler) http.Handler {
	cors := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers // Максимальное значение взято из гита chi cors
	})
	return cors
}
