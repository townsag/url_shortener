package main

import (
	"net"
	"net/http"
	"log"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"

	"townsag/url_shortener/src/middleware"
	"townsag/url_shortener/src/handlers"
)

func newServer(conn *pgx.Conn, rdb *redis.Client) http.Handler {
	mux := http.NewServeMux()
	handlers.AddRoutes(
		mux,
		conn,
		rdb,
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
	ctx := context.Background()
	var postgresConfig *dbConfig = getConfiguration()
	conn, err := createDBConnection(ctx, postgresConfig)
	if err != nil {
		log.Fatalf("failed to create a database connection: %s", err)
	}
	defer conn.Close(context.Background())

	var redisConfig *redisConfig = getRedisConfiguration()
	rdb, err := createRedisConnection(ctx, redisConfig)
	if err != nil {
		log.Fatalf("failed to create a redis connection: %s", err)
	}
	defer rdb.Close()

	srv := newServer(conn, rdb)
	httpServer := &http.Server{
		Addr: net.JoinHostPort("0.0.0.0", "8000"),
		Handler: srv,
	}

	log.Println("listening on port 8000")
	httpServer.ListenAndServe()
}