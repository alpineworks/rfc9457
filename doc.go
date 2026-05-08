// Package rfc9457 is a Go implementation of RFC 9457 "Problem Details for
// HTTP APIs".
//
// A problem detail describes a specific occurrence of an HTTP error in a
// machine-readable, structured format. Build one with [NewRFC9457] and the
// With* options, or with one of the status-named constructors such as
// [NotFound] or [BadRequest].
//
// *RFC9457 implements [http.Handler], [error], [encoding/json.Marshaler], and
// [encoding/json.Unmarshaler], so a problem can be served directly, returned
// from functions as a Go error, and round-tripped through encoding/json.
//
// https://datatracker.ietf.org/doc/html/rfc9457
package rfc9457
