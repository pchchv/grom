package grom

import (
	"net/http"
	"reflect"
)

// Request wraps net/http's Request and gocraf/web specific fields.
// In particular, PathParams is used to access captures params in your URL.
// A Request is sent to handlers on each request.
type Request struct {
	*http.Request
	// PathParams exists if you have wildcards in your URL that you need to capture.
	// Eg, /users/:id/tickets/:ticket_id and /users/1/tickets/33 would yield the map {id: "3", ticket_id: "33"}
	PathParams    map[string]string
	route         *route        // The actual route that got invoked.
	rootContext   reflect.Value // Root context. Set immediately.
	targetContext reflect.Value // The target context corresponding to the route. Not set until root middleware is done.
}
