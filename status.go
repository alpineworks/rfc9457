package rfc9457

import "net/http"

// problemForStatus is the shared constructor for status-named helpers. It
// seeds Status and Title (from http.StatusText), then applies caller options
// so they can override or extend.
func problemForStatus(status int, options []RFC9457Option) *RFC9457 {
	base := []RFC9457Option{
		WithStatus(status),
		WithTitle(http.StatusText(status)),
	}
	return NewRFC9457(append(base, options...)...)
}

// BadRequest builds a problem with status 400 and the matching standard title.
func BadRequest(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusBadRequest, options)
}

// Unauthorized builds a problem with status 401 and the matching standard title.
func Unauthorized(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusUnauthorized, options)
}

// Forbidden builds a problem with status 403 and the matching standard title.
func Forbidden(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusForbidden, options)
}

// NotFound builds a problem with status 404 and the matching standard title.
func NotFound(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusNotFound, options)
}

// MethodNotAllowed builds a problem with status 405 and the matching standard title.
func MethodNotAllowed(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusMethodNotAllowed, options)
}

// Conflict builds a problem with status 409 and the matching standard title.
func Conflict(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusConflict, options)
}

// Gone builds a problem with status 410 and the matching standard title.
func Gone(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusGone, options)
}

// UnprocessableEntity builds a problem with status 422 and the matching standard title.
func UnprocessableEntity(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusUnprocessableEntity, options)
}

// TooManyRequests builds a problem with status 429 and the matching standard title.
func TooManyRequests(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusTooManyRequests, options)
}

// InternalServerError builds a problem with status 500 and the matching standard title.
func InternalServerError(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusInternalServerError, options)
}

// NotImplemented builds a problem with status 501 and the matching standard title.
func NotImplemented(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusNotImplemented, options)
}

// BadGateway builds a problem with status 502 and the matching standard title.
func BadGateway(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusBadGateway, options)
}

// ServiceUnavailable builds a problem with status 503 and the matching standard title.
func ServiceUnavailable(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusServiceUnavailable, options)
}

// GatewayTimeout builds a problem with status 504 and the matching standard title.
func GatewayTimeout(options ...RFC9457Option) *RFC9457 {
	return problemForStatus(http.StatusGatewayTimeout, options)
}
