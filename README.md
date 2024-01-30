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
