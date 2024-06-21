package static

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// FileServer - настраивает обработчик на отдачу статических файлов.
func FileServer(r chi.Router, path string, root http.Handler) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, root)
		fs.ServeHTTP(w, r)
	})
}
