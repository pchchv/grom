package grom

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

var (
	// DefaultNotFoundResponse is the default text rendered when no route is found and no NotFound handlers are present.
	DefaultNotFoundResponse = "Not Found"
	// DefaultPanicResponse is the default text rendered when a panic occurs and no Error handlers are present.
	DefaultPanicResponse = "Application Error"
)

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

// If there's a panic in the root middleware (so that we don't have a route/target),
// then invoke the root handler or default.
// If there's a panic in other middleware, then invoke the target action's function.
// If there's a panic in the action handler, then invoke the target action's function.
func (rootRouter *Router) handlePanic(rw *appResponseWriter, req *Request, err interface{}) {
	var targetRouter *Router  // This will be set to the router we want to use the errorHandler on.
	var context reflect.Value // this is the context of the target router
	if req.route == nil {
		targetRouter = rootRouter
		context = req.rootContext
	} else {
		targetRouter = req.route.Router
		context = req.targetContext
		for !targetRouter.errorHandler.IsValid() && targetRouter.parent != nil {
			targetRouter = targetRouter.parent
			// Need to set context to the next context, UNLESS the context is the same type.
			curContextStruct := reflect.Indirect(context)
			if targetRouter.contextType != curContextStruct.Type() {
				context = curContextStruct.Field(0)
				if reflect.Indirect(context).Type() != targetRouter.contextType {
					panic("bug: shouldn't get here")
				}
			}
		}
	}

	if targetRouter.errorHandler.IsValid() {
		invoke(targetRouter.errorHandler, context, []reflect.Value{reflect.ValueOf(rw), reflect.ValueOf(req), reflect.ValueOf(err)})
	} else {
		http.Error(rw, DefaultPanicResponse, http.StatusInternalServerError)
	}

	const size = 4096
	stack := make([]byte, size)
	stack = stack[:runtime.Stack(stack, false)]

	PanicHandler.Panic(fmt.Sprint(req.URL), err, string(stack))
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
