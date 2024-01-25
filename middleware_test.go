package grom

import "fmt"

func (c *Context) A(w ResponseWriter, r *Request) {
	fmt.Fprintf(w, "context-A")
}

func (c *AdminContext) B(w ResponseWriter, r *Request) {
	fmt.Fprintf(w, "admin-B")
}

func (c *APIContext) C(w ResponseWriter, r *Request) {
	fmt.Fprintf(w, "api-C")
}

func (c *TicketsContext) D(w ResponseWriter, r *Request) {
	fmt.Fprintf(w, "tickets-D")
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

func (c *Context) mwGamma(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	fmt.Fprintf(w, "context-mw-Gamma ")
	next(w, r)
}

func (c *APIContext) mwDelta(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	fmt.Fprintf(w, "api-mw-Delta ")
	next(w, r)
}

func (c *AdminContext) mwEpsilon(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	fmt.Fprintf(w, "admin-mw-Epsilon ")
	next(w, r)
}

func (c *AdminContext) mwZeta(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	fmt.Fprintf(w, "admin-mw-Zeta ")
	next(w, r)
}

func (c *TicketsContext) mwEta(w ResponseWriter, r *Request, next NextMiddlewareFunc) {
	fmt.Fprintf(w, "tickets-mw-Eta ")
	next(w, r)
}
