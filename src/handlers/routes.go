package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"

	"townsag/url_shortener/src/db"
	// ^reference: https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html
	"townsag/url_shortener/src/util"
)

const ID_LENGTH int = 10

type createMappingRequestBody struct {
	LongUrl string `json:"long_url"`
}

type createMappingResponseBody struct {
	Msg string			`json:"message"`
	Status int			`json:"status"`
	ShortUrl *string	`json:"shortURL,omitempty"`
}


func createMappingHanlerFactory(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse the http post request body
		var body createMappingRequestBody
		err := util.DecodeJSONBody(w, r, body)
		if err != nil {
			var mr *util.MalformedRequest
			if errors.As(err, &mr) {
				// errors.As expects that the second argument is a pointer to an interface
				// this is why we have to us the address of mr. This is super odd right?
				http.Error(w, mr.Msg, mr.Status)
			} else {
				// TODO: log the error here with bound logger
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
		// write the long url to the database with retry
		queries := db.New(conn)
		var resultId string
		for i := 0; i < 3; i++ {
			tempResultId, err := util.RandomBase62(ID_LENGTH)
			if err != nil {
				// TODO: log that we were unable to generate a random number
				//		 what would even cause this?
				continue
			}
			params := db.InsertMappingParams{
				ID: tempResultId,
				LongUrl: body.LongUrl,
			}
			resultId, err = queries.InsertMapping(r.Context(), params)
			if err != nil || resultId == "" {
				// this means that either there was a database error or the
				// randomly generated result id is already assigned
				// TODO: log the error
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
		json.NewEncoder(w).Encode(response)
		// TODO: log the error from Encode
	}
}

func redirectToLongUrlHandlerFactory(conn *pgx.Conn) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// parse the short url from the path
		shortUrl := r.PathValue("shortUrlId")
		// TODO: add error handling for if the shortUrlId is not valid
		// read the relevant record from the database
		queries := db.New(conn)
		var record db.UrlMapping
		record, err := queries.SelectMapping(r.Context(), shortUrl)
		if err == pgx.ErrNoRows {
			http.Error(
				w, 
				fmt.Sprintf("could not find a mapping for shortUrlId: %s", shortUrl), 
				http.StatusNotFound,
			)
		}
		if err != nil {
			// TODO: get the logger from middleware
			// TODO: log the error here
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		// return a redirect to the long url associated with that short url
		http.Redirect(w, r, record.LongUrl, http.StatusFound)
	}
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func AddRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	// TODO: pass config / clients from main into the add routes function
	// HandleFunc under the hood creates and HandlerFunc object from the hello function and 
	// assignes that HandlerFunc object to the relevant pattern
	// The HandlerFunc type is just a function with an ServeHttp method defined on it that calls
	// the function
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("GET /{shortUrlId}", redirectToLongUrlHandlerFactory(conn))
	mux.HandleFunc("POST /mapping", createMappingHanlerFactory(conn))
}