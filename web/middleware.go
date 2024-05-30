package web

import (
	"github.com/go-chi/chi/v5/middleware"
)

func UseBasicAuth(r Router, username, password string) {
	if password != "" {
		if username == "" {
			username = "admin"
		}
		r.Use(middleware.BasicAuth("xlp", map[string]string{username: password}))
	}
}

var Recoverer = middleware.Recoverer
