package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"log/slog"

	"townsag/url_shortener/api/middleware"
)

func uiHandlersFactory(filesystem http.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logger *slog.Logger = middleware.GetLoggerFromContext(r.Context())
		path := r.URL.Path
		logger.Debug("received a request for a file", "path", path)
		// try if file exists at path, if not append .html
		_, err := filesystem.Open(path)
		if errors.Is(err, os.ErrNotExist) {
			path = fmt.Sprintf("%s.html", path)
		}
		r.URL.Path = path
		http.FileServer(filesystem).ServeHTTP(w, r)
	}
}