package middleware

import (
	"net/http"
	"github.com/google/uuid"
)

var requestIdHeader string = "X-Request-ID"

func RequestIdMiddleware(next http.Handler) http.Handler {
	// remember the pattern
	// middleware is a function that takes a http.Handler and returns another http Handler
	// This pattern uses closures to capture the next (input) http handler and wrap its
	// execution with something
	// http.HanderFunc is a type, casting this anonymous function to the HandlerFunc type
	// allows us to create a http.Handler instance from a function because http.HandlerFunc
	// implements the required method for it to conform to the http.Handler interface
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get(requestIdHeader)
		if requestId == "" {
			requestId = uuid.New().String()
			r.Header.Set(requestIdHeader, requestId)
		}
		next.ServeHTTP(w, r)
	})
}

/*
Still not sure if this is the right approach, I like the idea of encapsulating the
logic of fetching the id inside of the middleware package but this might be unecissary?
*/
func IdFromRequest(r *http.Request) (string) {
	requestId := r.Header.Get(requestIdHeader)
	return requestId
}