package datroute

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRouter(t *testing.T) {
	rs := Routes{
		Get("/", testHandlerFunc),
		Get("/profile/:name", testHandlerFunc),
		Get("/profile/:name", testHandlerFuncHelloWorld), // This should override the previous route.
		Post("/profile/:name", testHandlerFunc),
		Head("/profile/:name", testHandlerFunc),
		Put("/profile/:name", testHandlerFunc),
		Delete("/profile/:name", testHandlerFunc),
		Do("GET", "/profile/update", testHandlerFunc),
	}

	// Creating a new router manually.
	r := NewRouter()
	err := r.Handle(rs).Build()
	if err != nil {
		t.Errorf("Failed to build a handler. Error: %s.", err)
	}

	server := httptest.NewServer(r)
	defer server.Close()

	// Using a Build shortcut.
	h, err := rs.Build()
	if err != nil {
		t.Errorf("Failed to build a handler. Error: %s.", err)
	}

	server1 := httptest.NewServer(h)
	defer server1.Close()

	for _, v := range []struct {
		status                 int
		method, path, expected string
	}{
		{
			200, "GET", "/",
			fmt.Sprintf("method: GET, path: /, form: %v", url.Values{}),
		},
		{
			200, "GET", "/profile/john",
			fmt.Sprintf("Hello, world!\nmethod: GET, path: /profile/john, form: %v", url.Values{
				"name": {"john"},
			}),
		},
		{
			200, "POST", "/profile/jane",
			fmt.Sprintf("method: POST, path: /profile/jane, form: %v", url.Values{
				"name": {"jane"},
			}),
		},
		{
			200, "HEAD", "/profile/james", "",
		},
		{
			200, "PUT", "/profile/alice",
			fmt.Sprintf("method: PUT, path: /profile/alice, form: %v", url.Values{
				"name": {"alice"},
			}),
		},
		{
			200, "DELETE", "/profile/bob",
			fmt.Sprintf("method: DELETE, path: /profile/bob, form: %v", url.Values{
				"name": {"bob"},
			}),
		},
		{
			200, "GET", "/profile/update",
			fmt.Sprintf("method: GET, path: /profile/update, form: %v", url.Values{}),
		},
		{
			405, "POST", "/", "405 method not allowed\n",
		},
		{
			404, "POST", "/qwerty", "404 page not found\n",
		},
	} {
		for _, s := range []*httptest.Server{server, server1} {
			req, err := http.NewRequest(v.method, s.URL+v.path, nil)
			if err != nil {
				t.Errorf("Failed to create a new request. Error: %s.", err)
				continue
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("Cannot do a request. Error: %s.", err)
				continue
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("Did not manage to read a response body. Error: %s.", err)
			}
			actual := string(body)
			if res.StatusCode != v.status || actual != v.expected {
				t.Errorf(
					`%s "%s" => %#v %#v, expected %#v %#v.`,
					v.method, v.path, res.StatusCode, actual, v.status, v.expected,
				)
			}
		}
	}

}

func testHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "method: %s, path: %s, form: %v", r.Method, r.URL.Path, r.Form)
}

func testHandlerFuncHelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!\n")
	testHandlerFunc(w, r)
}
