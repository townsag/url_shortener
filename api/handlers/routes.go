
package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AddRoutes(
	mux *http.ServeMux, 
	conn *pgx.Conn, 
	rdb *redis.Client, 
	filesystem http.FileSystem,
	registry *prometheus.Registry,
) {
	// HandleFunc under the hood creates and HandlerFunc object from the hello function and 
	// assigns that HandlerFunc object to the relevant pattern
	// The HandlerFunc type is just a function with an ServeHttp method defined on it that calls
	// the function
	mux.Handle("GET /api/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	mux.HandleFunc("GET /api/healthy", healthyHandlerFactory(conn, rdb))
	mux.HandleFunc("GET /api/{shortUrlId}", redirectToLongUrlHandlerFactory(conn, rdb))
	mux.HandleFunc("POST /api/mapping", createMappingHandlerFactory(conn))
	mux.HandleFunc("/", uiHandlersFactory(filesystem))
}