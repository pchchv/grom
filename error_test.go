package grom

import (
	"fmt"
	"net/http"
	"testing"
)

func ErrorHandlerWithNoContext(w ResponseWriter, r *Request, err interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Contextless Error")
}

func TestNoErrorHandler(t *testing.T) {
	router := New(Context{})
	router.Get("/action", (*Context).ErrorAction)

	admin := router.Subrouter(AdminContext{}, "/admin")
	admin.Get("/action", (*AdminContext).ErrorAction)

	rw, req := newTestRequest("GET", "/action")
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "Application Error", http.StatusInternalServerError)

	rw, req = newTestRequest("GET", "/admin/action")
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "Application Error", http.StatusInternalServerError)
}

func TestHandlerOnRoot(t *testing.T) {
	router := New(Context{})
	router.Error((*Context).ErrorHandler)
	router.Get("/action", (*Context).ErrorAction)

	admin := router.Subrouter(AdminContext{}, "/admin")
	admin.Get("/action", (*AdminContext).ErrorAction)

	rw, req := newTestRequest("GET", "/action")
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "My Error", http.StatusInternalServerError)

	rw, req = newTestRequest("GET", "/admin/action")
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "My Error", http.StatusInternalServerError)
}
