package main

import (
	"errors"
	"log"
	"net/http"

	"alpineworks.io/rfc9457"
)

func main() {
	// Compose a problem from individual options.
	http.Handle("/test", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rfc9457.NewRFC9457(
			rfc9457.WithStatus(http.StatusTeapot),
			rfc9457.WithDetail("I'm a teapot"),
			rfc9457.WithTitle("Test endpoint hit - I'm a teapot"),
			rfc9457.WithInstance("/test"),
		).ServeHTTP(w, req)
	}))

	// Use a status-named constructor for common HTTP problems. Status and
	// title come from the helper; layer on detail/instance/extensions as
	// needed.
	http.Handle("/users/42", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rfc9457.NotFound(
			rfc9457.WithDetail("user 42 not found"),
			rfc9457.WithInstance(req.URL.Path),
			rfc9457.WithExtensions(rfc9457.NewExtension("user_id", 42)),
		).ServeHTTP(w, req)
	}))

	// Wrap an existing Go error.
	http.Handle("/wrap", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := errors.New("database connection failed")
		rfc9457.FromError(err, http.StatusInternalServerError).ServeHTTP(w, req)
	}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// curl localhost:8080/test
// {
//   "type": "about:blank",
//   "status": 418,
//   "title": "Test endpoint hit - I'm a teapot",
//   "detail": "I'm a teapot",
//   "instance": "/test"
// }
//
// curl localhost:8080/users/42
// {
//   "type": "about:blank",
//   "status": 404,
//   "title": "Not Found",
//   "detail": "user 42 not found",
//   "instance": "/users/42",
//   "user_id": 42
// }
