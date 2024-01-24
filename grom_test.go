package grom

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type nullPanicReporter struct{}

func (l nullPanicReporter) Panic(url string, err interface{}, stack string) {
	// no op
}

func init() {
	// This disables printing panics to stderr during testing, because that is very noisy,
	// and we purposefully test some panics.
	PanicHandler = nullPanicReporter{}
}

// Some default contexts and possible error handlers / actions
type Context struct{}

func (c *Context) ErrorMiddleware(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	var x, y int
	fmt.Fprintln(w, x/y)
}

func (c *Context) ErrorHandler(w ResponseWriter, r *Request, err interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "My Error")
}

func (c *Context) ErrorHandlerSecondary(w ResponseWriter, r *Request, err interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "My Secondary Error")
}

func (c *Context) ErrorAction(w ResponseWriter, r *Request) {
	var x, y int
	fmt.Fprintln(w, x/y)
}

type APIContext struct {
	*Context
}

// callerInfo returns the caller's caller info.
func callerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}
