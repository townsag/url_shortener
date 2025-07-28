package middleware

import (
	"log/slog"
	"os"
	"net/http"
	"context"
)

type contextKey string
const loggerKey contextKey = contextKey("logger")

func BuildLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(
		os.Stdout, 
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	))
	return logger
}

func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	// create a HandlerFunction instance from an anonymous function that closes
	// over the logger and the next middleware handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a bound logger with some request metadata
		boundLogger := logger.With(
			"method", r.Method,
			"path", r.URL.Path,
			requestIdHeader, IdFromRequest(r),
		)
		// add the logger to the request context
		ctx := context.WithValue(r.Context(), loggerKey, boundLogger)
		r = r.WithContext(ctx)
		// log some metadata about the request
		boundLogger.Debug("recieved request")
		// add call the next handler
		next.ServeHTTP(w, r)
	})
}

func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		return BuildLogger()
	}
	return logger
}