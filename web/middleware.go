package web

import (
	"github.com/go-chi/chi/v5/middleware"
)

// UseBasicAuth implements a simple middleware handler for adding basic http auth to a route.
func UseBasicAuth(r Router, username, password string) {
	if password != "" {
		if username == "" {
			username = "admin"
		}
		r.Use(middleware.BasicAuth("xlp", map[string]string{username: password}))
	}
}

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/go-chi/httplog middleware pkgs.
var Recoverer = middleware.Recoverer
