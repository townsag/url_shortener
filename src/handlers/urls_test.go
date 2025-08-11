package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
)

// test main calls the other tests in the handlers package testing suite
// use test main to make a package scoped postgres testcontainer
func TestMain(m *testing.M) {
	code := m.Run()

	cleanupPostgresContainer()
	cleanupRedisContainer()
	os.Exit(code)
}

func TestCreateMappingHappyPath(t *testing.T) {
	// get a connection to the postgres test container
	conn, err := setupPostgresContainer()
	if err != nil {
		t.Fatal(err)
	}
	// err = cleanupPostgresData()
	// if err != nil {
	// 	t.Fatalf("failed to restore the postgres database to the empty checkpoint %s", err)
	// }
	// create a create mapping handler
	handler := createMappingHandlerFactory(conn)
	// create a request for the create mapping route
	body := []byte(`{
		"longUrl": "https://google.com"
	}`)
	req, err := http.NewRequest("POST", "/mapping", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	// create a httptest Recorder instance. This implements the http.ResponseWriter
	// interface and allows us to record the response from the handler
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got: %v want %v", status, http.StatusOK)
		t.Fatalf("response body: %v", rr.Body)
	}

	var responseBody createMappingResponseBody
	decoder := json.NewDecoder(rr.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&responseBody)

	if err != nil {
		t.Fatalf("failed to decode response body with error: %v", err)
	}

	if responseBody.ShortUrl == nil {
		t.Fatalf(
			"failed to create a short url mapping, expected non nil value for ShortUrl field in response, got: %v", 
			responseBody.ShortUrl,
		)
	}
}

func TestCreateAndAccessMapping(t *testing.T) {
	conn, err := setupPostgresContainer()
	if err != nil {
		t.Fatal(err)
	}

	rdb, err := setupRedisContainer()
	if err != nil {
		t.Fatal(err)
	}
	
	testMux := http.NewServeMux()
	testMux.HandleFunc("GET /{shortUrlId}", redirectToLongUrlHandlerFactory(conn, rdb))
	testMux.HandleFunc("POST /mapping", createMappingHandlerFactory(conn))

	// for this test, assume that the create mapping call succeeds because failures of the
	// create mapping path will be caught by the other test
	body := []byte(`{"longUrl": "https://google.com"}`)
	req, err := http.NewRequest("POST", "/mapping", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testMux.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("create mapping returned incorrect response code: expected: %d, received: %d", http.StatusOK, status)
	}

	var responseBody createMappingResponseBody
	decoder := json.NewDecoder(rr.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&responseBody)

	if err != nil {
		t.Errorf("failed to decode create mapping response body with %v", err)
		t.Fatalf("response body: %v", rr.Body)
	}
	if responseBody.ShortUrl == nil {
		t.Fatal("failed to create a short url")
	}

	// call the redirect handler with the returned short url
	req, err = http.NewRequest("GET", fmt.Sprintf("/%s", *responseBody.ShortUrl), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	testMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("redirect to long url route returned incorrect status code: expected: %d, received: %d", 
		http.StatusFound, 
		status,
	)
		t.Fatalf("response body: %v", rr.Body)
	}
	if redirectLocation := rr.Result().Header.Get("Location"); redirectLocation != "https://google.com" {
		t.Fatalf("received unexpected redirect location: expected: https://google.com, received: %s", redirectLocation)
	}

	// verify that the long url is now in the cache
	// we use look aside caching so the long url is only written to the cache
	// on the read path
	longUrl, err := rdb.Get(context.Background(), *responseBody.ShortUrl).Result()
	if err != nil {
		if err == redis.Nil {
			t.Fatalf(
				"unable to find long url in redis cache for short url id %s, received: %v, expected: %s",
				*responseBody.ShortUrl,
				err,
				*responseBody.ShortUrl,
			)
		}
		t.Fatalf("error encountered when accessing redis cache: %v", err)
	}
	if longUrl != "https://google.com" {
		t.Fatalf(
			"retrieved a wrong value from the redis cache for a stored long url; expected: %s,  received: %s",
			"https://google.com",
			longUrl,
		)
	}
}

func TestAccessUnboundShortUrl(t *testing.T) {
	conn, err := setupPostgresContainer()
	if err != nil {
		t.Fatal(err)
	}
	
	rdb, err := setupRedisContainer()
	if err != nil {
		t.Fatal(err)
	}

	testMux := http.NewServeMux()
	testMux.HandleFunc("GET /{shortUrlId}", redirectToLongUrlHandlerFactory(conn, rdb))

	req, err := http.NewRequest("GET", "/12345678", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	testMux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Fatalf("handler returned incorrect status code: expected: %d, got: %d", http.StatusNotFound, status)
	}
}

func TestCreateInvalidMapping(t *testing.T) {
	// what makes a mapping invalid...?
	// should I be checking that the LongUrl is correctly formed or still up?
}

func TestAccessInvalidShortUrl(t *testing.T) {
	conn, err := setupPostgresContainer()
	if err != nil {
		t.Fatal(err)
	}

	rdb, err := setupRedisContainer()
	if err != nil {
		t.Fatal(err)
	}

	handler := redirectToLongUrlHandlerFactory(conn, rdb)

	req, err := http.NewRequest("GET", "/asdf", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned incorrect status code: expected: %d, got: %d", http.StatusBadRequest, status)
	}
}