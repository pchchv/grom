package grom

import "strings"

func (c *Context) OptionsHandler(rw ResponseWriter, req *Request, methods []string) {
	rw.Header().Add("Access-Control-Allow-Methods", strings.Join(methods, ", "))
	rw.Header().Add("Access-Control-Max-Age", "100")
}

func AccessControlMiddleware(rw ResponseWriter, req *Request, next NextMiddlewareFunc) {
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	next(rw, req)
}
