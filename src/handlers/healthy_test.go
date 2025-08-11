package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	// "context"
	// "github.com/jackc/pgx/v5"
)

func TestHealthySuccess(t *testing.T) {
	conn, err := setupPostgresContainer()
	if err != nil {
		t.Fatalf("failed to create a connection to postgres: %v", err)
	}

	dbr, err := setupRedisContainer()
	if err != nil {
		t.Fatalf("failed to create a connection to redis: %v", err)
	}

	testMux := http.NewServeMux()
	testMux.HandleFunc("GET /healthy", healthyHandlerFactory(conn, dbr))

	req, err := http.NewRequest("GET", "/healthy", nil)
	if err != nil {
		t.Fatalf("unable to make get request for healthy route %s", err)
	}
	rr := httptest.NewRecorder()
	testMux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("received wrong status code for healthy route, expected: %d, received: %d", http.StatusOK, rr.Code)
	}
}

type failingPinger struct{}
func (*failingPinger) Ping(ctx context.Context) error {
	return fmt.Errorf("failed to connect to database")
}

func TestHealthyFailPg(t *testing.T) {
	dbr, err := setupRedisContainer()
	if err != nil {
		t.Fatalf("failed to create a connection to redis: %v", err)
	}

	testMux := http.NewServeMux()
	testMux.HandleFunc("GET /healthy", healthyHandlerFactory(&failingPinger{}, dbr))

	req, err := http.NewRequest("GET", "/healthy", nil)
	if err != nil {
		t.Fatalf("unable to make a healthy request: %s", err)
	}
	rr := httptest.NewRecorder()
	testMux.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("received wring status code for healthy route, expected: %d, received: %d", http.StatusServiceUnavailable, rr.Code)
	}

}