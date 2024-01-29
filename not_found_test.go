package grom

import (
	"fmt"
	"net/http"
)

func (c *Context) HandlerWithContext(rw ResponseWriter, r *Request) {
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(rw, "My Not Found With Context")
}

func MyNotFoundHandler(rw ResponseWriter, r *Request) {
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(rw, "My Not Found")
}
