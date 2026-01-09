package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	webDir string
}

func NewHandler(webDir string) *Handler {
	return &Handler{
		webDir: webDir,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	// ТОЛЬКО статические файлы для Шага 1-2
	fileServer(r, "/", http.Dir(h.webDir))
}

// fileServer настраивает обработку статических файлов
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
