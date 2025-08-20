
package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func AddRoutes(mux *http.ServeMux, conn *pgx.Conn, rdb *redis.Client, filesystem http.FileSystem) {
	// HandleFunc under the hood creates and HandlerFunc object from the hello function and 
	// assigns that HandlerFunc object to the relevant pattern
	// The HandlerFunc type is just a function with an ServeHttp method defined on it that calls
	// the function
	mux.HandleFunc("GET /api/healthy", healthyHandlerFactory(conn, rdb))
	mux.HandleFunc("GET /api/{shortUrlId}", redirectToLongUrlHandlerFactory(conn, rdb))
	mux.HandleFunc("POST /api/mapping", createMappingHandlerFactory(conn))
	mux.HandleFunc("/", uiHandlersFactory(filesystem))
}