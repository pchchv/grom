package grom

import "net/http"

// Null response writer
type NullWriter struct{}

func (w *NullWriter) Header() http.Header {
	return nil
}

func (w *NullWriter) Write(data []byte) (n int, err error) {
	return len(data), nil
}

func (w *NullWriter) WriteHeader(statusCode int) {}
