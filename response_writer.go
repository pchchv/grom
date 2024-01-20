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

type appResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// Don't need this yet because we get it for free:
func (w *appResponseWriter) Write(data []byte) (n int, err error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	size, err := w.ResponseWriter.Write(data)
	w.size += size
	return size, err
}

func (w *appResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *appResponseWriter) Written() bool {
	return w.statusCode != 0
}

func (w *appResponseWriter) Size() int {
	return w.size
}
