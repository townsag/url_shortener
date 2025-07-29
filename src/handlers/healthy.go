package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type pinger interface {
	Ping(ctx context.Context) error
}

func healthyHandlerFactory(conn pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(r.Context(), time.Millisecond*500)
		if err := conn.Ping(ctx); err != nil {
			http.Error(
				w, 
				fmt.Sprintf("unable to connect to database: %s", err), 
				http.StatusServiceUnavailable,
			)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "healthy")
		}
	}
}