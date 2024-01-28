package grom

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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

func gromHandler(rw ResponseWriter, r *Request) {
	fmt.Fprintf(rw, "hello")
}

func gromRouterFor(namespaces []string, resources []string) http.Handler {
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

func BenchmarkGrom_Simple(b *testing.B) {
	router := New(BenchContext{})
	router.Get("/action", gromHandler)

	rw, req := testRequest("GET", "/action")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, req)
	}
}

func BenchmarkGrom_Middleware(b *testing.B) {
	router := New(BenchContext{})
	router.Middleware((*BenchContext).Middleware)
	router.Middleware((*BenchContext).Middleware)
	routerB := router.Subrouter(BenchContextB{}, "/b")
	routerB.Middleware((*BenchContextB).Middleware)
	routerB.Middleware((*BenchContextB).Middleware)
	routerC := routerB.Subrouter(BenchContextC{}, "/c")
	routerC.Middleware((*BenchContextC).Middleware)
	routerC.Middleware((*BenchContextC).Middleware)
	routerC.Get("/action", (*BenchContextC).Action)

	rw, req := testRequest("GET", "/b/c/action")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rw, req)
	}
}

func testRequest(method, path string) (*httptest.ResponseRecorder, *http.Request) {
	request, _ := http.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()

	return recorder, request
}

// Returns a routeset with N *resources per namespace*. so N=1 gives about 15 routes
func resourceSetup(N int) (namespaces []string, resources []string, requests []*http.Request) {
	namespaces = []string{"admin", "api", "site"}
	resources = []string{}

	for i := 0; i < N; i++ {
		sha1 := sha1.New()
		io.WriteString(sha1, fmt.Sprintf("%d", i))
		strResource := fmt.Sprintf("%x", sha1.Sum(nil))
		resources = append(resources, strResource)
	}

	for _, ns := range namespaces {
		for _, res := range resources {
			req, _ := http.NewRequest("GET", "/"+ns+"/"+res, nil)
			requests = append(requests, req)
			req, _ = http.NewRequest("POST", "/"+ns+"/"+res, nil)
			requests = append(requests, req)
			req, _ = http.NewRequest("GET", "/"+ns+"/"+res+"/3937", nil)
			requests = append(requests, req)
			req, _ = http.NewRequest("PUT", "/"+ns+"/"+res+"/3937", nil)
			requests = append(requests, req)
			req, _ = http.NewRequest("DELETE", "/"+ns+"/"+res+"/3937", nil)
			requests = append(requests, req)
		}
	}

	return
}

func benchmarkRoutesN(b *testing.B, N int, builder RouterBuilder) {
	namespaces, resources, requests := resourceSetup(N)
	router := builder(namespaces, resources)
	benchmarkRoutes(b, router, requests)
}

func benchmarkRoutes(b *testing.B, handler http.Handler, requests []*http.Request) {
	recorder := &NullWriter{}
	reqID := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if reqID >= len(requests) {
			reqID = 0
		}
		req := requests[reqID]
		handler.ServeHTTP(recorder, req)
		reqID++
	}
}
