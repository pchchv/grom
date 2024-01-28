package grom

import (
	"fmt"
	"net/http"
)

// Null response writer
type NullWriter struct{}

func (w *NullWriter) Header() http.Header {
	return nil
}

func (w *NullWriter) Write(data []byte) (n int, err error) {
	return len(data), nil
}

func (w *NullWriter) WriteHeader(statusCode int) {}

// Types used by any/all frameworks:
type RouterBuilder func(namespaces []string, resources []string) http.Handler

// Benchmarks for gocraft/web:
type BenchContext struct {
	MyField string
}

type BenchContextB struct {
	*BenchContext
}

type BenchContextC struct {
	*BenchContextB
}

func (c *BenchContext) Action(w ResponseWriter, r *Request) {
	fmt.Fprintf(w, "hello")
}

func (c *BenchContextB) Action(w ResponseWriter, r *Request) {
	fmt.Fprintf(w, c.MyField)
}

func (c *BenchContextC) Action(w ResponseWriter, r *Request) {
	fmt.Fprintf(w, "hello")
}

func (c *BenchContext) Middleware(rw ResponseWriter, r *Request, next NextMiddlewareFunc) {
	next(rw, r)
}

func (c *BenchContextB) Middleware(rw ResponseWriter, r *Request, next NextMiddlewareFunc) {
	next(rw, r)
}

func (c *BenchContextC) Middleware(rw ResponseWriter, r *Request, next NextMiddlewareFunc) {
	next(rw, r)
}
