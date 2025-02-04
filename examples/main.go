package main

import (
	"net/http"

	"alpineworks.io/rfc9457"
)

func main() {
	http.Handle("/test", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rfc9457.NewRFC9457(
			rfc9457.WithStatus(http.StatusTeapot),
			rfc9457.WithDetail("I'm a teapot"),
			rfc9457.WithTitle("Test endpoint hit - I'm a teapot"),
			rfc9457.WithInstance("/test"),
		).ServeHTTP(w, req)
	}))

	http.ListenAndServe(":8080", nil)
}

// {
// 	"detail": "I'm a teapot",
// 	"instance": "/test",
// 	"status": 418,
// 	"title": "Test endpoint hit - I'm a teapot",
// 	"type": "about:blank"
// }
