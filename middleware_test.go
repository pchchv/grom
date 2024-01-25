package grom

import "fmt"

func (c *Context) A(w ResponseWriter, r *Request) {
	fmt.Fprintf(w, "context-A")
}

func (c *Context) Z(w ResponseWriter, r *Request) {
	fmt.Fprintf(w, "context-Z")
}

func (c *Context) mwNoNext(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	fmt.Fprintf(w, "context-mw-NoNext ")
}

func (c *Context) mwAlpha(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	fmt.Fprintf(w, "context-mw-Alpha ")
	next(w, r)
}

func (c *Context) mwBeta(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	fmt.Fprintf(w, "context-mw-Beta ")
	next(w, r)
}
