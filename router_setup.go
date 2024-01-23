package grom

import (
	"reflect"
	"strings"
)

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

// New returns a new router with context type ctx.
// ctx should be a struct instance,
// whose purpose is to communicate type information.
// On each request,
// an instance of this context type will be automatically allocated and sent to handlers.
func New(ctx interface{}) *Router {
	validateContext(ctx, nil)
	r := &Router{}
	r.contextType = reflect.TypeOf(ctx)
	r.pathPrefix = "/"
	r.maxChildrenDepth = 1
	r.root = make(map[httpMethod]*pathNode)
	for _, method := range httpMethods {
		r.root[method] = newPathNode()
	}
	return r
}

// NewWithPrefix returns a new router (see New) but each route will have an implicit prefix.
// For instance, with pathPrefix = "/api/v2",
// all routes under this router will begin with "/api/v2".
func NewWithPrefix(ctx interface{}, pathPrefix string) *Router {
	r := New(ctx)
	r.pathPrefix = pathPrefix
	return r
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

// Since it's easy to pass the wrong method to a middleware/handler route,
// and since the user can't rely on static type checking since we use reflection,
// lets be super helpful about what they did and what they need to do.
// Arguments:
//   - vfn is the failed method
//   - addingType is for "You are adding {addingType} to a router...". E.g. "middleware" or "a handler" or "an error handler"
//   - yourType is for "Your {yourType} function can have...". Eg, "middleware" or "handler" or "error handler"
//   - args is like "rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc"
//   - NOTE: args can be calculated if you pass in each type. BUT, it doesn't have example argument name, so it has less copy/paste value.
func instructiveMessage(vfn reflect.Value, addingType string, yourType string, args string, ctxType reflect.Type) string {
	// Get context type without package.
	ctxString := ctxType.String()
	splitted := strings.Split(ctxString, ".")
	if len(splitted) <= 1 {
		ctxString = splitted[0]
	} else {
		ctxString = splitted[1]
	}

	str := "\n" + strings.Repeat("*", 120) + "\n"
	str += "* You are adding " + addingType + " to a router with context type '" + ctxString + "'\n"
	str += "*\n*\n"
	str += "* Your " + yourType + " function can have one of these signatures:\n"
	str += "*\n"
	str += "* // If you don't need context:\n"
	str += "* func YourFunctionName(" + args + ")\n"
	str += "*\n"
	str += "* // If you want your " + yourType + " to accept a context:\n"
	str += "* func (c *" + ctxString + ") YourFunctionName(" + args + ")  // or,\n"
	str += "* func YourFunctionName(c *" + ctxString + ", " + args + ")\n"
	str += "*\n"
	str += "* Unfortunately, your function has this signature: " + vfn.Type().String() + "\n"
	str += "*\n"
	str += strings.Repeat("*", 120) + "\n"
	return str
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

// Panics unless fn is a proper handler wrt ctxType
// eg, func(ctx *ctxType, writer, request)
func validateHandler(vfn reflect.Value, ctxType reflect.Type) {
	var req *Request
	var resp func() ResponseWriter
	if !isValidHandler(vfn, ctxType, reflect.TypeOf(resp).Out(0), reflect.TypeOf(req)) {
		panic(instructiveMessage(vfn, "a handler", "handler", "rw web.ResponseWriter, req *web.Request", ctxType))
	}
}

func validateErrorHandler(vfn reflect.Value, ctxType reflect.Type) {
	var req *Request
	var resp func() ResponseWriter
	if !isValidHandler(vfn, ctxType, reflect.TypeOf(resp).Out(0), reflect.TypeOf(req), emptyInterfaceType) {
		panic(instructiveMessage(vfn, "an error handler", "error handler", "rw web.ResponseWriter, req *web.Request, err interface{}", ctxType))
	}
}

func validateNotFoundHandler(vfn reflect.Value, ctxType reflect.Type) {
	var req *Request
	var resp func() ResponseWriter
	if !isValidHandler(vfn, ctxType, reflect.TypeOf(resp).Out(0), reflect.TypeOf(req)) {
		panic(instructiveMessage(vfn, "a 'not found' handler", "not found handler", "rw web.ResponseWriter, req *web.Request", ctxType))
	}
}

func validateOptionsHandler(vfn reflect.Value, ctxType reflect.Type) {
	var req *Request
	var resp func() ResponseWriter
	var methods []string
	if !isValidHandler(vfn, ctxType, reflect.TypeOf(resp).Out(0), reflect.TypeOf(req), reflect.TypeOf(methods)) {
		panic(instructiveMessage(vfn, "an 'options' handler", "options handler", "rw web.ResponseWriter, req *web.Request, methods []string", ctxType))
	}
}

func validateMiddleware(vfn reflect.Value, ctxType reflect.Type) {
	var req *Request
	var resp func() ResponseWriter
	var n NextMiddlewareFunc
	if !isValidHandler(vfn, ctxType, reflect.TypeOf(resp).Out(0), reflect.TypeOf(req), reflect.TypeOf(n)) {
		panic(instructiveMessage(vfn, "middleware", "middleware", "rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc", ctxType))
	}
}
