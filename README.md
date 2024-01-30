# grom [![GoDoc](https://godoc.org/github.com/pchchv/grom?status.png)](https://godoc.org/github.com/pchchv/grom)


**grom** is a Go mux and middleware package.

## Features
* **Ultrafast and scalable**. The added latency is 3-9µs per request. Routing performance is O(log(N)) in number of routes.
* **Your own contexts**. Easily pass information between middleware and handler with strong static typing.
* **Easy and powerful routing**. Intercept path variables. Validate path segments with regexps.
* **Middleware**. Middleware can express almost any web-layer feature.
* **Nested routers, contexts, and middleware**. Your application has an API, an admin area, and a view for logging out. Each view needs different contexts and different middleware. 
* **Use the go net/http package**. Start your server with http.ListenAndServe(), and work directly with http.ResponseWriter and http.Request.
* **Minimal**. The gocraft/web core is lightweight and minimal. Add additional functionality with built-in middleware or write your own middleware.

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
