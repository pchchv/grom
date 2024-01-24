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

func (mw *middlewareHandler) invoke(ctx reflect.Value, rw ResponseWriter, req *Request, next NextMiddlewareFunc) {
	if mw.Generic {
		mw.GenericMiddleware(rw, req, next)
	} else {
		mw.DynamicMiddleware.Call([]reflect.Value{ctx, reflect.ValueOf(rw), reflect.ValueOf(req), reflect.ValueOf(next)})
	}
}

// routersFor returns [root router, child router, ..., leaf route's router]
// given the route and the target router.
// Uses memory in routers to store this information.
func routersFor(route *route, routers []*Router) []*Router {
	routers = routers[:0]
	curRouter := route.Router
	for curRouter != nil {
		routers = append(routers, curRouter)
		curRouter = curRouter.parent
	}

	// Reverse the slice
	s, e := 0, len(routers)-1
	for s < e {
		routers[s], routers[e] = routers[e], routers[s]
		s++
		e--
	}
	return routers
}

// contexts is initially filled with a single context for the root
// routers is [root, child, ..., leaf] with at least 1 element
// Returns [ctx for root, ... ctx for leaf]
// NOTE: if two routers have the same contextType, then they'll share the exact same context.
func contextsFor(contexts []reflect.Value, routers []*Router) []reflect.Value {
	routersLen := len(routers)
	for i := 1; i < routersLen; i++ {
		var ctx reflect.Value
		if routers[i].contextType == routers[i-1].contextType {
			ctx = contexts[i-1]
		} else {
			ctx = reflect.New(routers[i].contextType)
			// set the first field to the parent
			f := reflect.Indirect(ctx).Field(0)
			f.Set(contexts[i-1])
		}
		contexts = append(contexts, ctx)
	}
	return contexts
}

func invoke(handler reflect.Value, ctx reflect.Value, values []reflect.Value) {
	numIn := handler.Type().NumIn()
	if numIn == len(values) {
		handler.Call(values)
	} else {
		values = append([]reflect.Value{ctx}, values...)
		handler.Call(values)
	}
}

func calculateRoute(rootRouter *Router, req *Request) (*route, map[string]string) {
	var leaf *pathLeaf
	var wildcardMap map[string]string
	method := httpMethod(req.Method)
	tree, ok := rootRouter.root[method]
	if ok {
		leaf, wildcardMap = tree.Match(req.URL.Path)
	}

	// If no match and this is a HEAD, route on GET.
	if leaf == nil && method == httpMethodHead {
		tree, ok := rootRouter.root[httpMethodGet]
		if ok {
			leaf, wildcardMap = tree.Match(req.URL.Path)
		}
	}

	if leaf == nil {
		return nil, nil
	}
	return leaf.route, wildcardMap
}
