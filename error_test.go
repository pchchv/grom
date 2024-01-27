package grom

import (
	"fmt"
	"net/http"
)

func ErrorHandlerWithNoContext(w ResponseWriter, r *Request, err interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Contextless Error")
}
