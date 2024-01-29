package grom

import "strings"

func (c *Context) OptionsHandler(rw ResponseWriter, req *Request, methods []string) {
	rw.Header().Add("Access-Control-Allow-Methods", strings.Join(methods, ", "))
	rw.Header().Add("Access-Control-Max-Age", "100")
}
