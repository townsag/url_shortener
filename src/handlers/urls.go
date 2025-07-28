package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/jackc/pgx/v5"

	"townsag/url_shortener/src/db"
	"townsag/url_shortener/src/middleware"

	// ^reference: https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html
	"townsag/url_shortener/src/util"
)

const ID_LENGTH int = 8

type createMappingRequestBody struct {
	LongUrl string `json:"longUrl"`
}

type createMappingResponseBody struct {
	Msg string			`json:"message"`
	Status int			`json:"status"`
	ShortUrl *string	`json:"shortUrl,omitempty"`
}


func createMappingHandlerFactory(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logger *slog.Logger = middleware.GetLoggerFromContext(r.Context())
		// parse the http post request body
		var body createMappingRequestBody
		err := util.DecodeJSONBody(w, r, &body)
		if err != nil {
			var mr *util.MalformedRequest
			if errors.As(err, &mr) {
				logger.Warn("client error encountered when validating request body", "error", err)
				// errors.As expects that the second argument is a pointer to an interface
				// this is why we have to us the address of mr. This is super odd right?
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(mr.Status)
				json.NewEncoder(w).Encode(mr)
				return
			} else {
				logger.Error("server error encountered when decoding request body", "error", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(&createMappingResponseBody{
					Msg: http.StatusText(http.StatusInternalServerError),
					Status: http.StatusInternalServerError,
				})
				return
			}
		}
		// write the long url to the database with retry
		queries := db.New(conn)
		var resultId string
		for i := range 3 {
			tempResultId, err := util.RandomBase62(ID_LENGTH)
			if err != nil {
				logger.Error("failed to generate a short url", "error", err)
				continue
			}
			params := db.InsertMappingParams{
				ID: tempResultId,
				LongUrl: body.LongUrl,
			}
			resultId, err = queries.InsertMapping(r.Context(), params)
			if err != nil {
				logger.Error("database error encountered when writing new long url", "error", err)
				continue
			}
			if resultId == "" {
				logger.Warn("tried to insert duplicate short url", "attempt", i)
				continue
			}
			break
		}
		var response createMappingResponseBody
		if resultId == "" {
			response = createMappingResponseBody{
				Msg: "failed to create short url because of internal server error",
				Status: http.StatusInternalServerError,
			}
		} else {
			response = createMappingResponseBody{
				Msg: "successfully created short url",
				Status: http.StatusOK,
				ShortUrl: &resultId,
			}
		}
		// return the generated short url
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		// TODO: log the error from Encode
	}
}

type redirectToLongUrlResponseBody struct {
	Msg string	`json:"message"`
	Status int	`json:"status"`
}

func isValidShortUrlId(id string) bool {
	r := regexp.MustCompile(fmt.Sprintf("^[a-zA-Z0-9]{%d}$", ID_LENGTH))
	return r.MatchString(id)
}

func redirectToLongUrlHandlerFactory(conn *pgx.Conn) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var logger *slog.Logger = middleware.GetLoggerFromContext(r.Context())
		// parse the short url from the path
		shortUrl := r.PathValue("shortUrlId")
		// error handling for if the shortUrlId is not valid
		if !isValidShortUrlId(shortUrl) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&redirectToLongUrlResponseBody{
				Msg: fmt.Sprintf("received invalid url mapping id: %s, must be %d characters long and include only [a-zA-Z0-9]", shortUrl, ID_LENGTH),
				Status: http.StatusBadRequest,
			})
		}
		// read the relevant record from the database
		queries := db.New(conn)
		var record db.UrlMapping
		record, err := queries.SelectMapping(r.Context(), shortUrl)
		if err == pgx.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&redirectToLongUrlResponseBody{
				Msg: fmt.Sprintf("could not find a mapping for shortUrlId: %s", shortUrl),
				Status: http.StatusNotFound,
			})
			return
		}
		if err != nil {
			logger.Error(
				"database error encountered when querying for long url", 
				"error", err, 
				"shortUrl", shortUrl,
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&redirectToLongUrlResponseBody{
				Msg: http.StatusText(http.StatusInternalServerError),
				Status: http.StatusInternalServerError,
			})
			return
		}
		// return a redirect to the long url associated with that short url
		http.Redirect(w, r, record.LongUrl, http.StatusFound)
	}
}