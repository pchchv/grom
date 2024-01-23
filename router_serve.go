package grom

import "reflect"

type middlewareClosure struct {
	appResponseWriter
	Request
	Routers                []*Router
	Contexts               []reflect.Value
	currentMiddlewareIndex int
	currentRouterIndex     int
	currentMiddlewareLen   int
	RootRouter             *Router
	Next                   NextMiddlewareFunc
}
