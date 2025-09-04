
package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func AddRoutes(
	mux *http.ServeMux, 
	conn *pgx.Conn, 
	rdb *redis.Client, 
	filesystem http.FileSystem,
) {
	// The HandlerFunc type is just a function with an ServeHttp method defined on it that calls
	// the function

	mux.Handle("GET /api/healthy", otelhttp.WithRouteTag("GET /api/healthy", healthyHandlerFactory(conn, rdb)))
	mux.Handle("GET /api/{shortUrlId}", otelhttp.WithRouteTag("GET /api/{shortUrlId}", redirectToLongUrlHandlerFactory(conn, rdb)))
	mux.Handle("POST /api/mapping", otelhttp.WithRouteTag("POST /api/mapping", createMappingHandlerFactory(conn)))
	mux.Handle("/", otelhttp.WithRouteTag("/", uiHandlersFactory(filesystem)))
}