package grom

const (
	httpMethodGet     = httpMethod("GET")
	httpMethodPost    = httpMethod("POST")
	httpMethodPut     = httpMethod("PUT")
	httpMethodDelete  = httpMethod("DELETE")
	httpMethodPatch   = httpMethod("PATCH")
	httpMethodHead    = httpMethod("HEAD")
	httpMethodOptions = httpMethod("OPTIONS")
)

type httpMethod string

// NextMiddlewareFunc are functions passed into your middleware. To advance the middleware, call the function.
// You should usually pass the existing ResponseWriter and *Request into the next middlware, but you can
// chose to swap them if you want to modify values or capture things written to the ResponseWriter.
type NextMiddlewareFunc func(ResponseWriter, *Request)
