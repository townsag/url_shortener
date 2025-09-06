package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func AddRoutes(
	mux *http.ServeMux, 
	pool *pgxpool.Pool, 
	rdb *redis.Client, 
	filesystem http.FileSystem,
) {
	// The HandlerFunc type is just a function with an ServeHttp method defined on it that calls
	// the function

	mux.Handle("GET /api/healthy", otelhttp.WithRouteTag("GET /api/healthy", healthyHandlerFactory(pool, rdb)))
	mux.Handle("GET /api/{shortUrlId}", otelhttp.WithRouteTag("GET /api/{shortUrlId}", redirectToLongUrlHandlerFactory(pool, rdb)))
	mux.Handle("POST /api/mapping", otelhttp.WithRouteTag("POST /api/mapping", createMappingHandlerFactory(pool)))
	mux.Handle("/", otelhttp.WithRouteTag("/", uiHandlersFactory(filesystem)))
}