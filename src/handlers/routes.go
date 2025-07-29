
package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

func AddRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	// HandleFunc under the hood creates and HandlerFunc object from the hello function and 
	// assigns that HandlerFunc object to the relevant pattern
	// The HandlerFunc type is just a function with an ServeHttp method defined on it that calls
	// the function
	mux.HandleFunc("GET /healthy", healthyHandlerFactory(conn))
	mux.HandleFunc("GET /{shortUrlId}", redirectToLongUrlHandlerFactory(conn))
	mux.HandleFunc("POST /mapping", createMappingHandlerFactory(conn))
}