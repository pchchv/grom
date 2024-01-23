package grom

import "reflect"

const (
	httpMethodGet     = httpMethod("GET")
	httpMethodPost    = httpMethod("POST")
	httpMethodPut     = httpMethod("PUT")
	httpMethodDelete  = httpMethod("DELETE")
	httpMethodPatch   = httpMethod("PATCH")
	httpMethodHead    = httpMethod("HEAD")
	httpMethodOptions = httpMethod("OPTIONS")
)

var (
	httpMethods        = []httpMethod{httpMethodGet, httpMethodPost, httpMethodPut, httpMethodDelete, httpMethodPatch, httpMethodHead, httpMethodOptions}
	emptyInterfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
)

type httpMethod string

// NextMiddlewareFunc are functions passed into your middleware. To advance the middleware, call the function.
// You should usually pass the existing ResponseWriter and *Request into the next middlware, but you can
// chose to swap them if you want to modify values or capture things written to the ResponseWriter.
type NextMiddlewareFunc func(ResponseWriter, *Request)

// GenericHandler are handlers that don't have or need a context. If your handler doesn't need a context,
// you can use this signature to get a small performance boost.
type GenericHandler func(ResponseWriter, *Request)

// GenericMiddleware are middleware that doesn't have or need a context. General purpose middleware, such as
// static file serving, has this signature. If your middlware doesn't need a context, you can use this
// signature to get a small performance boost.
type GenericMiddleware func(ResponseWriter, *Request, NextMiddlewareFunc)

type actionHandler struct {
	Generic        bool
	DynamicHandler reflect.Value
	GenericHandler GenericHandler
}

type route struct {
	Router  *Router
	Method  httpMethod
	Path    string
	Handler *actionHandler
}

type middlewareHandler struct {
	Generic           bool
	DynamicMiddleware reflect.Value
	GenericMiddleware GenericMiddleware
}

// Router implements net/http's Handler interface and is what you attach middleware, routes/handlers, and subrouters to.
type Router struct {
	// Hierarchy:
	parent           *Router // nil if root router.
	children         []*Router
	maxChildrenDepth int
	// For each request we'll create one of these objects
	contextType reflect.Type
	// e.g. "/" or "/admin".
	// Any routes added to this router will be prefixed with this.
	pathPrefix string
	// Routeset contents:
	middleware []*middlewareHandler
	routes     []*route
	// The root pathnode is the same for a tree of Routers
	root map[httpMethod]*pathNode
	// This can can be set on any router.
	// The target's ErrorHandler will be invoked if it exists.
	errorHandler reflect.Value
	// This can only be set on the root handler, since by virtue of not finding a route, we don't have a target.
	// (That being said, in the future we could investigate namespace matches)
	notFoundHandler reflect.Value
	// This can only be set on the root handler, since by virtue of not finding a route, we don't have a target.
	optionsHandler reflect.Value
}

// Ensures vfn is a function, that optionally takes a *ctxType as the first argument,
// followed by the specified types.
// Handlers have no return value.
// Returns true if valid, false otherwise.
func isValidHandler(vfn reflect.Value, ctxType reflect.Type, types ...reflect.Type) bool {
	fnType := vfn.Type()
	if fnType.Kind() != reflect.Func {
		return false
	}

	typesStartIdx := 0
	typesLen := len(types)
	numIn := fnType.NumIn()
	numOut := fnType.NumOut()
	if numOut != 0 {
		return false
	}

	if numIn == typesLen {
		// No context
	} else if numIn == (typesLen + 1) {
		// context, types
		firstArgType := fnType.In(0)
		if firstArgType != reflect.PtrTo(ctxType) && firstArgType != emptyInterfaceType {
			return false
		}
		typesStartIdx = 1
	} else {
		return false
	}

	for _, typeArg := range types {
		if fnType.In(typesStartIdx) != typeArg {
			return false
		}
		typesStartIdx++
	}
	return true
}

// Panics unless validation is correct
func validateContext(ctx interface{}, parentCtxType reflect.Type) {
	ctxType := reflect.TypeOf(ctx)
	if ctxType.Kind() != reflect.Struct {
		panic("web: Context needs to be a struct type")
	}

	if parentCtxType != nil && parentCtxType != ctxType {
		if ctxType.NumField() == 0 {
			panic("web: Context needs to have first field be a pointer to parent context")
		}

		fldType := ctxType.Field(0).Type
		// Ensure fld is a pointer to parentCtxType
		if fldType != reflect.PtrTo(parentCtxType) {
			panic("web: Context needs to have first field be a pointer to parent context")
		}
	}
}
