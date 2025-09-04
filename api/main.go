package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"townsag/url_shortener/api/handlers"
	"townsag/url_shortener/api/middleware"
)

//go:embed all:build
var files embed.FS

func newServer(conn *pgx.Conn, rdb *redis.Client, filesystem http.FileSystem) http.Handler {
	mux := http.NewServeMux()
	handlers.AddRoutes(
		mux,
		conn,
		rdb,
		filesystem,
	)

	root_logger := middleware.BuildLogger()

	var handler http.Handler = mux
	// applying the middleware in this order means that the request id middleware
	// will execute and then the logging middleware
	handler = middleware.LoggingMiddleware(root_logger, handler)
	handler = middleware.RequestIdMiddleware(handler)
	handler = otelhttp.NewHandler(
		handler,
		"url-shortener",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		}),
	)
	return handler
}

func main() {
	ctx := context.Background()

	// bootstrap the OTEL SDK
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		log.Fatalf("failed to bootstrap OTEL SDK: %s", err)
	}
	defer otelShutdown(context.Background())
	// TODO: do something with that error^

	// create a connection to the postgres database server
	var postgresConfig *dbConfig = getConfiguration()
	conn, err := createDBConnection(ctx, postgresConfig)
	if err != nil {
		log.Fatalf("failed to create a database connection: %s", err)
	}
	defer conn.Close(context.Background())

	// create a connection to the redis server
	var redisConfig *redisConfig = getRedisConfiguration()
	rdb, err := createRedisConnection(ctx, redisConfig)
	if err != nil {
		log.Fatalf("failed to create a redis connection: %s", err)
	}
	defer rdb.Close()

	// create a filesystem object
	fsys, err := fs.Sub(files, "build")
	if err != nil {
		log.Fatalf("unable to access the filesystem: %s", err)
	}
	filesystem := http.FS(fsys)

	// build the server with its routes
	srv := newServer(conn, rdb, filesystem)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", "8000"),
		Handler: srv,
	}

	log.Println("listening on port 8000")
	httpServer.ListenAndServe()
}
