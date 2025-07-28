
package handlers

import (
	"net/http"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func AddRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	// HandleFunc under the hood creates and HandlerFunc object from the hello function and 
	// assignes that HandlerFunc object to the relevant pattern
	// The HandlerFunc type is just a function with an ServeHttp method defined on it that calls
	// the function
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("GET /{shortUrlId}", redirectToLongUrlHandlerFactory(conn))
	mux.HandleFunc("POST /mapping", createMappingHandlerFactory(conn))
}