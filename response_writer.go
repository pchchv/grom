package grom

import (
	"context"
	"net/http"
)

// ResponseWriter includes net/http's ResponseWriter and adds a StatusCode() method to obtain the written status code.
// A ResponseWriter is sent to handlers on each request.
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	http.Hijacker
	context.Context
	// StatusCode returns the written status code, or 0 if none has been written yet.
	StatusCode() int
	// Written returns whether the header has been written yet.
	Written() bool
	// Size returns the size in bytes of the body written so far.
	Size() int
}
