package datroute

import (
	"net/http"
)

// MethodNotAllowed replies to the request with an HTTP 405 method not allowed
// error. If you want to use your own MethodNotAllowed handler, please override
// this variable.
var MethodNotAllowed = func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}

// NotFound replies to the request with an HTTP 404 not found error.
// NotFound is called when unknown HTTP method or a handler not found.
// If you want to use the your own NotFound handler, please overwrite this variable.
var NotFound = func(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
