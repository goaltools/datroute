// Package datroute is a double array trie based router that uses
// Form field of request to pass parameters to handler functions.
// As a result, full compatability with the standard library is achieved.
// The downside is a memory consumption and speed are a bit worse
// than if it would pass a third Params argument to handler functions.
// The router is based on denco package.
//
// A sample of its usage is below:
//
//	package main
//
//	import (
//		"log"
//		"net/http"
//
//		r "github.com/goaltools/datroute"
//	)
//
//	func main() {
//		router := r.NewRouter()
//		err := router.Handle(r.Routes{
//			r.Get("/profiles/:username", Profiles.ShowUserHandleFunc),
//			r.Delete("/profiles/:username", Profiles.DeleteUserHandleFunc),
//		}).Build()
//		if err != nil {
//			log.Fatal(err)
//		}
//		log.Fatal(http.ListenAndServe(":8080", router))
//	}
package datroute

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/naoina/denco"
)

// Router represents a multiplexer for HTTP requests.
type Router struct {
	data    *denco.Router  // data stores denco router.
	indexes map[string]int // indexes are used to simplify the search of records we need.
	records []denco.Record // records is a list of handlers expected by denco router.
}

// Routes is an alias of []Route.
type Routes []*Route

// Route is used to store information about HTTP request's handler
// including a list of allowed methods and pattern.
type Route struct {
	Handlers *Dict  // HTTP request method -> handler pairs.
	Pattern  string // Pattern is a routing path for handler.
}

// NewRouter allocates and returns a new multiplexer.
func NewRouter() *Router {
	return &Router{
		indexes: map[string]int{},
	}
}

// Get is an short form of Route("GET", pattern, handler).
func Get(pattern string, handler http.HandlerFunc) *Route {
	return Do("GET", pattern, handler)
}

// Post is a short form of Route("POST", pattern, handler).
func Post(pattern string, handler http.HandlerFunc) *Route {
	return Do("POST", pattern, handler)
}

// Put is a short form of Route("PUT", pattern, handler).
func Put(pattern string, handler http.HandlerFunc) *Route {
	return Do("PUT", pattern, handler)
}

// Head is a short form of Route("HEAD", pattern, handler).
func Head(pattern string, handler http.HandlerFunc) *Route {
	return Do("HEAD", pattern, handler)
}

// Delete is a short form of Route("DELETE", pattern, handler).
func Delete(pattern string, handler http.HandlerFunc) *Route {
	return Do("DELETE", pattern, handler)
}

// ServeHTTP is used to implement http.Handler interface.
// It dispatches the request to the handler whose pattern
// most closely matches the request URL.
func (t *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, _ := t.Handler(r)
	h.ServeHTTP(w, r)
}

// Handle registers handlers for given patterns.
// If a handler already exists for pattern, it will be overridden.
// If it exists but with another method, a new method will be added.
func (t *Router) Handle(routes Routes) *Router {
	for i := range routes {
		// Check whether we have already had such route.
		index, ok := t.indexes[routes[i].Pattern]

		// If we haven't, add the route.
		if !ok {
			// Save pattern's index to simplify its search
			// in next iteration.
			t.indexes[routes[i].Pattern] = len(t.records)

			// Add the route to the slice.
			t.records = append(t.records, denco.NewRecord(routes[i].Pattern, routes[i]))
			continue
		}

		// Otherwise, just add new HTTP methods to the existing route.
		r := t.records[index].Value.(*Route)
		r.Handlers.Join(routes[i].Handlers)
	}
	return t
}

// Build compiles registered routes. Routes that are added after building will not
// be handled. A new call to build will be required.
func (t *Router) Build() error {
	t.data = denco.New()
	return t.data.Build(t.records)
}

// Build is a shortcut get an http handler for the routes.
func (t Routes) Build() (http.Handler, error) {
	r := NewRouter()
	err := r.Handle(t).Build()
	return r, err
}

// Do allocates and returns a Route struct.
func Do(method, pattern string, handler http.HandlerFunc) *Route {
	hs := NewDict()
	hs.Set(strings.ToUpper(method), &handler)
	return &Route{
		Handlers: hs,
		Pattern:  pattern,
	}
}

// Handler returns the handler to use for the given request, consulting r.Method
// and r.URL.Path. It always returns a non-nil handler. If there is no registered handler
// that applies to the request, Handler returns a “page not found” handler and empty pattern.
// If there is a registered handler but requested method is not allowed,
// "method not allowed" and a pattern are returned.
func (t *Router) Handler(r *http.Request) (handler http.Handler, pattern string) {
	// Make sure we have a handler for this request.
	obj, params, found := t.data.Lookup(r.URL.Path)
	if !found {
		return http.HandlerFunc(NotFound), ""
	}

	// Check whether requested method is allowed.
	route := obj.(*Route)
	handler, i := route.Handlers.Get(r.Method)
	if i == -1 {
		return http.HandlerFunc(MethodNotAllowed), route.Pattern
	}

	// Add parameters of request to request.Form and return a handler.
	if len(params) > 0 {
		r.Form = make(url.Values, len(params))
		for i := range params {
			r.Form[params[i].Name] = []string{params[i].Value}
		}
	}
	return handler, route.Pattern
}
