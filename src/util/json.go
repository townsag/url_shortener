package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	// ^ errors in golang are passed as values instead of raise like exceptions
	// an error in golang can wrap another error
	// errors.Is or errors.As can be used to inspect the type of the returned error
)

// reference: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

const ONE_MB int = 1048576

// create a struct type for malformed request that satisfies the Error interface
type MalformedRequest struct {
	Msg string
	Status int
}

func (m *MalformedRequest) Error() string {
	return m.Msg
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	// verify that the request does include json body
	if r.Header.Get("Content-Type") != "application/json" {
		// return a malformed request object
		return &MalformedRequest{
			Msg: "Content-Type header must be application/json",
			Status: http.StatusUnsupportedMediaType,
		}
	}

	// cap the size of the request body to one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(ONE_MB))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(destination)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var maxBytesError *http.MaxBytesError
		// switch with no issue is like a series of if, else if statements
		switch {
			// "An error matches target if the error's concrete value is assignable to the value pointed to by target"
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains a syntax error at position: %d", syntaxError.Offset)
			return &MalformedRequest{Msg: msg, Status: http.StatusBadRequest}
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "request body contains badly formatted json"
			return &MalformedRequest{Msg: msg, Status: http.StatusBadRequest}
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("request body contains and invalid value for the field %q at position %d", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &MalformedRequest{Msg: msg, Status: http.StatusBadRequest}
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("request body contains unknown field: %s", fieldName)
			return &MalformedRequest{Msg: msg, Status: http.StatusBadRequest}
		case errors.Is(err, io.EOF):
			msg := "request body must not be empty"
			return &MalformedRequest{Msg: msg, Status: http.StatusBadRequest}
		case errors.As(err, &maxBytesError):
			msg := fmt.Sprintf("request body must not exceed %d bytes", ONE_MB)
			return &MalformedRequest{Msg: msg, Status: http.StatusRequestEntityTooLarge}
		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})
	// ^create an instance of the Any type
	if !errors.Is(err, io.EOF) {
		msg := "request body must not have more than one JSON object"
		return &MalformedRequest{Msg: msg, Status: http.StatusBadRequest}
	}

	return nil
}