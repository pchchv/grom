package grom

import (
	"os"
	"strings"
	"testing"
)

func routerSetupBody() string {
	fileBytes, _ := os.ReadFile(testFilename())
	return string(fileBytes)
}

func testFilename() string {
	return "router_setup.go"
}

func TestStaticMiddleware(t *testing.T) {
	currentRoot, _ := os.Getwd()

	router := New(Context{})
	router.Middleware(StaticMiddleware(currentRoot))
	router.Get("/action", (*Context).A)

	// Make sure we can still hit actions:
	rw, req := newTestRequest("GET", "/action")
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "context-A", 200)

	rw, req = newTestRequest("GET", "/"+testFilename())
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, strings.TrimSpace(routerSetupBody()), 200)

	rw, req = newTestRequest("POST", "/"+testFilename())
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "Not Found", 404)
}

func TestStaticMiddlewareOptionPrefix(t *testing.T) {
	currentRoot, _ := os.Getwd()

	router := New(Context{})
	router.Middleware(StaticMiddleware(currentRoot, StaticOption{Prefix: "/public"}))
	router.Get("/action", (*Context).A)

	rw, req := newTestRequest("GET", "/action")
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "context-A", 200)

	rw, req = newTestRequest("GET", "/"+testFilename())
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "Not Found", 404)

	rw, req = newTestRequest("GET", "/public/"+testFilename())
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, strings.TrimSpace(routerSetupBody()), 200)
}
