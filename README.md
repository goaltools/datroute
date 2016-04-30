# DATRoute
DATRoute is a double array trie based router that uses
`Form` field of request to pass parameters to handler functions
for compatibility with the standard library.
It is based on [`denco`](https://github.com/naoina/denco) router.

[![GoDoc](https://godoc.org/github.com/goaltools/datroute?status.svg)](https://godoc.org/github.com/goaltools/datroute)
[![Build Status](https://travis-ci.org/goaltools/datroute.svg?branch=master)](https://travis-ci.org/goaltools/datroute)
[![Coverage](https://codecov.io/github/goaltools/datroute/coverage.svg?branch=master)](https://codecov.io/github/goaltools/datroute?branch=master)
[![Go Report Card](http://goreportcard.com/badge/goaltools/datroute?t=3)](http:/goreportcard.com/report/goaltools/datroute)

### Installation
*Use `-u` ("update") flag to make sure the latest version of package is installed.*
```bash
go get -u github.com/goaltools/datroute
```

### Usage
```go
package main

import (
	"log"
	"net/http"

	r "github.com/goaltools/datroute"
)

func main() {
	router := r.NewRouter()
	err := router.Handle(r.Routes{
		r.Get("/profiles/:username", ShowUserHandleFunc),
		r.Delete("/profiles/:username", DeleteUserHandleFunc),
	}).Build()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8080", router))
}
```

### License
Distributed under the BSD 2-clause "Simplified" License unless otherwise noted.
