package handlers
import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req * http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func AddRoutes(mux *http.ServeMux) {
	// TODO: pass config / clients from main into the add routes function
	// HandleFunc under the hood creates and HandlerFunc object from the hello function and 
	// assignes that HandlerFunc object to the relevant pattern
	// The HandlerFunc type is just a function with an ServeHttp method defined on it that calls
	// the function
	mux.HandleFunc("/hello", hello)
}