# grom [![GoDoc](https://godoc.org/github.com/pchchv/grom?status.png)](https://godoc.org/github.com/pchchv/grom)


**grom** is a Go mux and middleware package.

## Features
* **Ultrafast and scalable**. The added latency is 3-9µs per request. Routing performance is O(log(N)) in number of routes.
* **Your own contexts**. Easily pass information between middleware and handler with strong static typing.
* **Easy and powerful routing**. Intercept path variables. Validate path segments with regexps.
* **Middleware**. Middleware can express almost any web-layer feature.
* **Nested routers, contexts, and middleware**. Your application has an API, an admin area, and a view for logging out. Each view needs different contexts and different middleware. 
* **Use the go net/http package**. Start your server with http.ListenAndServe(), and work directly with http.ResponseWriter and http.Request.
* **Minimal**. The grom core is lightweight and minimal. Add additional functionality with built-in middleware or write your own middleware.

## Performance
Performance is paramount.

For minimalistic applications, the added latency is about 3 µs. For more complex applications (6 middleware functions, 3 levels of contexts, 150+ routes) this increases to 10µs.

One of the key decisions, is the choice of routing algorithm. Most libraries use a simple O(N) iteration over all routes to find a match. This is fine if you only have a few routes, but it starts to fail as the size of the application increases. Here is a tree-based router whose complexity grows by O(log(N)).

## Usage

```bash
go get github.com/pchchv/grom
```

Add a file ```server.go```

```go
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pchchv/grom"
)

type Context struct {
	HelloCount int
}

func (c *Context) SetHelloCount(rw grom.ResponseWriter, req *grom.Request, next grom.NextMiddlewareFunc) {
	c.HelloCount = 3
	next(rw, req)
}

func (c *Context) SayHello(rw grom.ResponseWriter, req *grom.Request) {
	fmt.Fprint(rw, strings.Repeat("Hello ", c.HelloCount), "World!")
}

func main() {
	router := grom.New(Context{}). // Create your router
					Middleware(grom.LoggerMiddleware).     // Use some included middleware
					Middleware(grom.ShowErrorsMiddleware). // ...
					Middleware((*Context).SetHelloCount).  // Your own middleware!
					Get("/", (*Context).SayHello)          // Add a route
	http.ListenAndServe("localhost:3000", router) // Start the server!
}
```

Run the server. It will be available on ```localhost:3000```:

```bash
go run server.go
```

## Application Structure

### Making your router
The first thing you need to do is make a new router. Routers serve requests and execute middleware.

```go
router := grom.New(YourContext{})
```

### Your context
YourContext{} - is any structure you want. For example:

```go
type YourContext struct {
    User *User // Assumes you've defined a User type as well.
}
```

Your context can be empty or it can have various fields in it. The fields can be whatever you want - it's your type! When a new request comes into the router, we'll allocate an instance of this struct and pass it to your middleware and handlers. This allows, for instance, a SetUser middleware to set a User field that can be read in the handlers.

### Routes and handlers
Once you have your router, you can add routes to it. Standard HTTP verbs are supported.

```go
router := grom.New(YourContext{})
router.Get("/users", (*YourContext).UsersList)
router.Post("/users", (*YourContext).UsersCreate)
router.Put("/users/:id", (*YourContext).UsersUpdate)
router.Delete("/users/:id", (*YourContext).UsersDelete)
router.Patch("/users/:id", (*YourContext).UsersUpdate)
router.Get("/", (*YourContext).Root)
```

``(*YourContext).Root`` is a method expression. It allows your handlers to look like this:

```go
func (c *YourContext) Root(rw grom.ResponseWriter, req *grom.Request) {
    if c.User != nil {
        fmt.Fprint(rw, "Hello,", c.User.Name)
	} else {
		fmt.Fprint(rw, "Hello, anonymous person")
	}
}
```

All method expressions do is return a function that accepts the type as the first argument. So your handler can also look like this:

```go
func Root(c *YourContext, rw grom.ResponseWriter, req *grom.Request) {}
```

Of course, if you don't need a context for a particular action, you can also do that:

```go
func Root(rw grom.ResponseWriter, req *grom.Request) {}
```

Note that handlers always need to accept two input parameters: grom.ResponseWriter, and *grom.Request, both of which wrap the standard http.ResponseWriter and *http.Request, respectively.

### Middleware
You can add middleware to a router:

```go
router := grom.New(YourContext{})
router.Middleware((*YourContext).UserRequired)
// add routes, more middleware
```

This is what a middleware handler looks like:

```go
func (c *YourContext) UserRequired(rw grom.ResponseWriter, r *grom.Request, next grom.NextMiddlewareFunc) {
	user := userFromSession(r)  // Pretend like this is defined. It reads a session cookie and returns a *User or nil.
	if user != nil {
		c.User = user
		next(rw, r)
	} else {
		rw.Header().Set("Location", "/")
		rw.WriteHeader(http.StatusMovedPermanently)
		// do NOT call next()
	}
}
```

Some things to note about the above example:
*  We set fields in the context for future middleware / handlers to use.
*  We can call next(), or not. Not calling next() effectively stops the middleware stack.

Of course, generic middleware without contexts is supported:

```go
func GenericMiddleware(rw grom.ResponseWriter, r *grom.Request, next grom.NextMiddlewareFunc) {
	// ...
}
```

### Nested routers
Nested routers allow you to run different middleware and use different contexts for different parts of the application. Some common scenarios are:
* You want to run AdminRequired middleware on all Admin routes, but not on API routes. Your context needs a CurrentAdmin field.
* You want to run the OAuth middleware on API routes. Your context needs the AccessToken field.
* You want to run session middleware on ALL of your routes. Your context needs the Session field.

Let's implement this. Your contexts will look like the following:

```go
type Context struct {
	Session map[string]string
}

type AdminContext struct {
	*Context
	CurrentAdmin *User
}

type ApiContext struct {
	*Context
	AccessToken string
}
```

Note that we embed a pointer to the parent context in each subcontext. This is required.

Now that we have our contexts, let's create our routers:

```go
rootRouter := grom.New(Context{})
rootRouter.Middleware((*Context).LoadSession)

apiRouter := rootRouter.Subrouter(ApiContext{}, "/api")
apiRouter.Middleware((*ApiContext).OAuth)
apiRouter.Get("/tickets", (*ApiContext).TicketsIndex)

adminRouter := rootRouter.Subrouter(AdminContext{}, "/admin")
adminRouter.Middleware((*AdminContext).AdminRequired)

// Given the path namesapce for this router is "/admin", the full path of this route is "/admin/reports"
adminRouter.Get("/reports", (*AdminContext).Reports)
```

Note that each time we make a subrouter, we need to supply the context as well as a path namespace. The context CAN be the same as the parent context, and the namespace CAN just be "/" for no namespace.

### Request lifecycle
The following is a detailed account of the request lifecycle:

1.  A request comes in. Yay! (follow along in ```router_serve.go``` if you'd like)
2.  Wrap the default Go http.ResponseWriter and http.Request in a grom.ResponseWriter and grom.Request, respectively (via structure embedding).
3.  Allocate a new root context. This context is passed into your root middleware.
4.  Execute middleware on the root router. We do this before we find a route!
5.  After all of the root router's middleware is executed, we'll run a 'virtual' routing middleware that determines the target route.
    *  If the there's no route found, we'll execute the NotFound handler if supplied. Otherwise, we'll write a 404 response and start unwinding the root middlware.
6.  Now that we have a target route, we can allocate the context tree of the target router.
7.  Start executing middleware on the nested middleware leading up to the final router/route.
8.  After all middleware is executed, we'll run another 'virtual' middleware that invokes the final handler corresponding to the target route.
9.  Unwind all middleware calls (if there's any code after next() in the middleware, obviously that's going to run at some point).

### Capturing path params; regexp conditions
You can capture path variables like this:

```go
router.Get("/suggestions/:suggestion_id/comments/:comment_id")
```

In your handler, you can access them like this:

```go
func (c *YourContext) Root(rw grom.ResponseWriter, req *grom.Request) {
	fmt.Fprint(rw, "Suggestion ID:", req.PathParams["suggestion_id"])
	fmt.Fprint(rw, "Comment ID:", req.PathParams["comment_id"])
}
```

You can also validate the format of your path params with a regexp. For instance, to ensure the 'ids' start with a digit:

```go
router.Get("/suggestions/:suggestion_id:\\d.*/comments/:comment_id:\\d.*")
```

You can match any route past a certain point like this:

```go
router.Get("/suggestions/:suggestion_id/comments/:comment_id/:*")
```

The path params will contain a “*” member with the rest of your path.  It is illegal to add any more paths past the “*” path param, as it’s meant to match every path afterwards, in all cases.

For Example:
    /suggestions/123/comments/321/foo/879/bar/834

Elicits path params:
    * “suggestion_id”: 123,
    * “comment_id”: 321,
    * “*”: “foo/879/bar/834”


One thing you CANNOT currently do is use regexps outside of a path segment. For instance, optional path segments are not supported - you would have to define multiple routes that both point to the same handler. This design decision was made to enable efficient routing.

### Not Found handlers
If a route isn't found, by default we'll return a 404 status and render the text "Not Found".

You can supply a custom NotFound handler on your root router:

```go
router.NotFound((*Context).NotFound)
```

Your handler can optionally accept a pointer to the root context. NotFound handlers look like this:

```go
func (c *Context) NotFound(rw grom.ResponseWriter, r *grom.Request) {
	rw.WriteHeader(http.StatusNotFound) // You probably want to return 404. But you can also redirect or do whatever you want.
	fmt.Fprintf(rw, "My Not Found")     // Render you own HTML or something!
}
```

### OPTIONS handlers
If an [OPTIONS request](https://en.wikipedia.org/wiki/Cross-origin_resource_sharing#Preflight_example) is made and routes with other methods are found for the requested path, then by default we'll return an empty response with an appropriate `Access-Control-Allow-Methods` header.

You can supply a custom OPTIONS handler on your root router:

```go
router.OptionsHandler((*Context).OptionsHandler)
```

Your handler can optionally accept a pointer to the root context. OPTIONS handlers look like this:

```go
func (c *Context) OptionsHandler(rw grom.ResponseWriter, r *grom.Request, methods []string) {
	rw.Header().Add("Access-Control-Allow-Methods", strings.Join(methods, ", "))
	rw.Header().Add("Access-Control-Allow-Origin", "*")
}
```

### Error handlers
By default, if there's a panic in middleware or a handler, we'll return a 500 status and render the text "Application Error".

If you use the included middleware ```grom.ShowErrorsMiddleware```, a panic will result in a pretty backtrace being rendered in HTML. This is great for development.

You can also supply a custom Error handler on any router (not just the root router):

```go
router.Error((*Context).Error)
```

Your handler can optionally accept a pointer to its corresponding context. Error handlers look like this:

```go
func (c *Context) Error(rw grom.ResponseWriter, r *grom.Request, err interface{}) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "Error", err)
}
```

### Included middleware
We ship with three basic pieces of middleware: a logger, an exception printer, and a static file server. To use them:

```go
router := grom.New(Context{})
router.Middleware(grom.LoggerMiddleware).
	Middleware(grom.ShowErrorsMiddleware)

// The static middleware serves files. Examples:
// "GET /" will serve an index file at pwd/public/index.html
// "GET /robots.txt" will serve the file at pwd/public/robots.txt
// "GET /images/foo.gif" will serve the file at pwd/public/images/foo.gif
currentRoot, _ := os.Getwd()
router.Middleware(grom.StaticMiddleware(path.Join(currentRoot, "public"), grom.StaticOption{IndexFile: "index.html"}))
```

NOTE: You might not want to use grom.ShowErrorsMiddleware in production. You can easily do something like this:
```go
router := grom.New(Context{})
router.Middleware(grom.LoggerMiddleware)
if MyEnvironment == "development" {
	router.Middleware(grom.ShowErrorsMiddleware)
}
// ...
```

### Starting your server
Since grom.Router implements http.Handler (eg, ServeHTTP(ResponseWriter, *Request)), you can easily plug it in to the standard Go http machinery:

```go
router := grom.New(Context{})
// ... Add routes and such.
http.ListenAndServe("localhost:8080", router)
```

### Rendering responses
So now you routed a request to a handler. You have a grom.ResponseWriter (http.ResponseWriter) and grom.Request (http.Request). Now what?

```go
// You can print to the ResponseWriter!
fmt.Fprintf(rw, "<html>I'm a web page!</html>")
```

This is currently where the implementation of this library stops. I recommend you read the documentation of [net/http](http://golang.org/pkg/net/http/).

## Toolkit

* [work](https://github.com/pchchv/work) - Process background jobs in Go.
* [grom](https://github.com/pchchv/grom) - Go Router + Middleware. Your Contexts.
* [health](https://github.com/pchchv/health) -  Instrument your web apps with logging and metrics.