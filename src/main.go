package main

import (
	"fmt"
	"net"
	"net/http"
	"townsag/url_shortener/src/middleware"
	"townsag/url_shortener/src/handlers"
)

func newServer() http.Handler {
	// TODO: add dependencies here like a database connection pool
	mux := http.NewServeMux()
	// TODO: pass the database connection pool into add routes
	handlers.AddRoutes(
		mux,
	)

	root_logger := middleware.BuildLogger()

	var handler http.Handler = mux
	// applying the middleware in this order means that the request id middleware
	// will execute and then the logging middleware
	handler = middleware.LoggingMiddleware(root_logger, handler)
	handler = middleware.RequestIdMiddleware(handler)
	return handler
}

func main() {
	srv := newServer()
	httpServer := &http.Server{
		Addr: net.JoinHostPort("0.0.0.0", "8000"),
		Handler: srv,
	}

	fmt.Println("listening on port 8000")
	httpServer.ListenAndServe()
}