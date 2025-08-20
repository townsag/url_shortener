package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type redisClient interface {
	Ping(ctx context.Context) *redis.StatusCmd
}

func healthyHandlerFactory(conn pinger, dbr redisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*500)
		defer cancel()
		if err := conn.Ping(ctx); err != nil {
			http.Error(
				w, 
				fmt.Sprintf("unable to connect to database: %s", err), 
				http.StatusServiceUnavailable,
			)
			return
		} else if err := dbr.Ping(ctx).Err(); err != nil {
			http.Error(
				w,
				fmt.Sprintf("unable to connect to redis cache: %s", err),
				http.StatusServiceUnavailable,
			)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "healthy")
		}
	}
}