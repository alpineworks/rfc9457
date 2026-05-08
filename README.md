<h1 align="center">
  <img src=".github/images/rfc9457.png" alt="rfc9457" width="500">
</h1>
<h2 align="center">
    rfc9457 provides an implementation of the RFC "Problem Details for HTTP APIs" in Go
</h2>

<div align="center">

[RFC 9457][rfc-link]

[![Alpineworks][alpineworks-badge]][for-the-badge-link] [![Made With Go][made-with-go-badge]][for-the-badge-link]

</div>

---

## Usage
See [examples/](examples/) for extra implementation information for this library.

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
	"type": "about:blank",
	"status": 418,
	"title": "Test endpoint hit - I'm a teapot",
	"detail": "I'm a teapot",
	"instance": "/test"
}
```
</details>

<details>
<summary>Custom problem type</summary>

`Type` defaults to `about:blank`. Use `WithType` to point at a documented problem-type URI:

```go
rfc9457.NewRFC9457(
	rfc9457.WithType("https://example.com/problems/insufficient-funds"),
	rfc9457.WithStatus(http.StatusForbidden),
	rfc9457.WithTitle("You do not have enough credit"),
	rfc9457.WithDetail("Your current balance is 30, but charge was 50"),
	rfc9457.WithInstance("/account/12345/transactions/abc"),
).ServeHTTP(w, req)
```

Per RFC 9457, when decoding a problem with a missing `type`, `FromJSON` populates `Type` with `about:blank`.
</details>

<details>
<summary>Status-named constructors</summary>

For common HTTP problems, the status-named helpers seed `Status` and a default `Title` from `http.StatusText`. Layer on detail/instance/extensions as needed:

```go
rfc9457.NotFound(
	rfc9457.WithDetail("user 42 not found"),
	rfc9457.WithInstance("/users/42"),
	rfc9457.WithExtensions(rfc9457.NewExtension("user_id", 42)),
).ServeHTTP(w, req)
```

Available helpers: `BadRequest`, `Unauthorized`, `Forbidden`, `NotFound`, `MethodNotAllowed`, `Conflict`, `Gone`, `UnprocessableEntity`, `TooManyRequests`, `InternalServerError`, `NotImplemented`, `BadGateway`, `ServiceUnavailable`, `GatewayTimeout`.
</details>

<details>
<summary>Extensions</summary>

Extension members flatten into the top-level JSON object on encode and are pulled back into `Extensions` on decode:

```go
rfc9457.NewRFC9457(
	rfc9457.WithStatus(http.StatusBadRequest),
	rfc9457.WithTitle("Validation failed"),
	rfc9457.WithExtensions(
		rfc9457.NewExtension("trace_id", "abc-123"),
		rfc9457.NewExtension("invalid_params", []map[string]string{
			{"name": "email", "reason": "must be a valid email"},
		}),
	),
).ServeHTTP(w, req)
```

`WithExtensions` is additive across calls — each invocation merges into the existing map (later keys overwrite earlier ones), so option helpers can compose cleanly.

Extension keys that collide with reserved members (`type`, `status`, `title`, `detail`, `instance`) cause `ToJSON` / `MarshalJSON` to return `ErrExtensionKeyCollision`.
</details>

<details>
<summary>Wrapping a Go error</summary>

`FromError` builds a problem with `Title` from `http.StatusText(status)` and `Detail` from `err.Error()`:

```go
if err := doWork(); err != nil {
	rfc9457.FromError(err, http.StatusInternalServerError).ServeHTTP(w, req)
	return
}
```

`*RFC9457` also implements the `error` interface, so a problem can be returned and inspected as a Go error.
</details>

<details>
<summary>Encoding and decoding</summary>

```go
// Encode
problem := rfc9457.BadRequest(rfc9457.WithDetail("missing field 'email'"))
body, err := problem.ToJSON() // []byte

// Decode
parsed, err := rfc9457.FromJSON(body)
```

`*RFC9457` also implements `json.Marshaler` and `json.Unmarshaler`, so it works directly with `encoding/json`.

Sentinel errors returned by this package: `ErrUnableToMarshalJSON`, `ErrUnableToUnmarshalJSON`, `ErrExtensionKeyCollision`. Use `errors.Is` to check them.
</details>

<details>
<summary>HTTP serving behavior</summary>

`*RFC9457` implements `http.Handler`. `ServeHTTP`:

- Sets `Content-Type` to `application/problem+json` (also exposed as `rfc9457.ProblemContentType`).
- Writes `r.Status` as the response code, defaulting to `500 Internal Server Error` when `Status` is unset.
- If serialization fails (for example, an extension key collides with a reserved field), writes a minimal fallback problem body with status 500 instead of leaving the response empty.
</details>

<!-- Reference Variables -->

<!-- Badges -->
[alpineworks-badge]: .github/images/alpine-works.svg
[made-with-go-badge]: .github/images/made-with-go.svg

<!-- Links -->
[blank-reference-link]: #
[for-the-badge-link]: https://forthebadge.com
[rfc-link]: https://www.rfc-editor.org/rfc/rfc9457.html
