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

// Benchmarks for pchchv/web:
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

func webHandler(rw ResponseWriter, r *Request) {
	fmt.Fprintf(rw, "hello")
}

func webRouterFor(namespaces []string, resources []string) http.Handler {
	router := New(BenchContext{})
	for _, ns := range namespaces {
		subrouter := router.Subrouter(BenchContext{}, "/"+ns)
		for _, res := range resources {
			subrouter.Get("/"+res, (*BenchContext).Action)
			subrouter.Post("/"+res, (*BenchContext).Action)
			subrouter.Get("/"+res+"/:id", (*BenchContext).Action)
			subrouter.Put("/"+res+"/:id", (*BenchContext).Action)
			subrouter.Delete("/"+res+"/:id", (*BenchContext).Action)
		}
	}
	return router
}
