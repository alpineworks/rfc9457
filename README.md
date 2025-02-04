<h1 align="center">
  <img src=".github/images/rfc9457.png" alt="rfc9457" width="500">
</h1>
<h2 align="center">
    rfc9457 provides an implementation of the RFC "Problem Details for HTTP APIs" in Go
</h2>

<div align="center">

[RFC 9457][rfc-link]    q

[![Alpineworks][alpineworks-badge]][for-the-badge-link] [![Made With Go][made-with-go-badge]][for-the-badge-link]

</div>

---

## Usage
See [examples/main.go](examples/main.go) for information in using the library.

<details open>
<summary>Suggested Implementation</summary>

```go
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
```

## Output
`curl localhost:8080/test`
```json
{
	"detail": "I'm a teapot",
	"instance": "/test",
	"status": 418,
	"title": "Test endpoint hit - I'm a teapot",
	"type": "about:blank"
}
```
</details>

<!-- Reference Variables -->

<!-- Badges -->
[alpineworks-badge]: .github/images/alpine-works.svg
[made-with-go-badge]: .github/images/made-with-go.svg

<!-- Links -->
[blank-reference-link]: #
[for-the-badge-link]: https://forthebadge.com
[rfc-link]: https://www.rfc-editor.org/rfc/rfc9457.html