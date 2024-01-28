package grom

import (
	"os"
	"path/filepath"
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

// TestStaticMiddlewareIndex will create an assets folder with one nested subfolder.
// Each folder will have an index.html file.
func TestStaticMiddlewareOptionIndex(t *testing.T) {
	// Create two temporary folders:
	dirName, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err.Error())
	}
	nestedDirName, err := os.MkdirTemp(dirName, "")
	if err != nil {
		panic(err.Error())
	}

	// Get the last path segment of the nestedDirName:
	_, nestedDirSegment := filepath.Split(nestedDirName)

	// Create first index file
	indexFilename := filepath.Join(dirName, "index.html")
	err = os.WriteFile(indexFilename, []byte("index1"), os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	defer os.Remove(indexFilename)

	// Create second index file
	nestedIndexFilename := filepath.Join(nestedDirName, "index.html")
	err = os.WriteFile(nestedIndexFilename, []byte("index2"), os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	defer os.Remove(nestedIndexFilename)

	// Make router. Static middleware rooted at first temp dir
	router := New(Context{})
	router.Middleware(StaticMiddleware(dirName, StaticOption{IndexFile: "index.html"}))
	router.Get("/action", (*Context).A)

	// Getting a root index:
	rw, req := newTestRequest("GET", "/")
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "index1", 200)

	// Nested dir
	rw, req = newTestRequest("GET", "/"+nestedDirSegment)
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "index2", 200)

	// Nested dir with trailing slash:
	rw, req = newTestRequest("GET", "/"+nestedDirSegment+"/")
	router.ServeHTTP(rw, req)
	assertResponse(t, rw, "index2", 200)
}
