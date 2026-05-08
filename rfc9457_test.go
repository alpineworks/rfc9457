package rfc9457_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"alpineworks.io/rfc9457"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRFC9457ToJSON(t *testing.T) {
	tests := []struct {
		Name         string
		Options      []rfc9457.RFC9457Option
		ExpectedJSON string
	}{
		{
			Name: "simple",
			Options: []rfc9457.RFC9457Option{
				rfc9457.WithTitle("failed to run query"),
			},
			ExpectedJSON: `{"type":"about:blank","title":"failed to run query"}`,
		},
		{
			Name: "with extensions",
			Options: []rfc9457.RFC9457Option{
				rfc9457.WithTitle("failed to run query"),
				rfc9457.WithExtensions(rfc9457.NewExtension("key", "value")),
			},
			ExpectedJSON: `{"type":"about:blank","title":"failed to run query","key":"value"}`,
		},
		{
			Name: "with status",
			Options: []rfc9457.RFC9457Option{
				rfc9457.WithTitle("failed to run query"),
				rfc9457.WithStatus(http.StatusInternalServerError),
			},
			ExpectedJSON: `{"type":"about:blank","title":"failed to run query","status":500}`,
		},
		{
			Name: "with instance",
			Options: []rfc9457.RFC9457Option{
				rfc9457.WithTitle("failed to run query"),
				rfc9457.WithInstance("/users/42"),
			},
			ExpectedJSON: `{"type":"about:blank","title":"failed to run query","instance":"/users/42"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			problem := rfc9457.NewRFC9457(tc.Options...)

			out, err := problem.ToJSON()
			require.NoError(t, err)

			assert.JSONEq(t, tc.ExpectedJSON, string(out))
		})
	}
}

func TestWithExtensionsMerges(t *testing.T) {
	problem := rfc9457.NewRFC9457(
		rfc9457.WithTitle("t"),
		rfc9457.WithExtensions(rfc9457.NewExtension("a", 1)),
		rfc9457.WithExtensions(rfc9457.NewExtension("b", 2)),
	)

	out, err := problem.ToJSON()
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"about:blank","title":"t","a":1,"b":2}`, string(out))
}

func TestExtensionKeyCollisionRejected(t *testing.T) {
	problem := rfc9457.NewRFC9457(
		rfc9457.WithTitle("t"),
		rfc9457.WithExtensions(rfc9457.NewExtension("status", 999)),
	)

	_, err := problem.ToJSON()
	require.Error(t, err)
	assert.ErrorIs(t, err, rfc9457.ErrExtensionKeyCollision)
}

func TestJSONToRFC9457(t *testing.T) {
	tests := []struct {
		Name            string
		JSON            string
		ExpectedRFC9457 *rfc9457.RFC9457
	}{
		{
			Name: "simple",
			JSON: `{"type":"about:blank","title":"failed to run query"}`,
			ExpectedRFC9457: &rfc9457.RFC9457{
				Type:  "about:blank",
				Title: "failed to run query",
			},
		},
		{
			Name: "with extensions",
			JSON: `{"type":"about:blank","title":"failed to run query","key":"value"}`,
			ExpectedRFC9457: &rfc9457.RFC9457{
				Type:  "about:blank",
				Title: "failed to run query",
				Extensions: map[string]any{
					"key": "value",
				},
			},
		},
		{
			Name: "with status",
			JSON: `{"type":"about:blank","title":"failed to run query","status":500}`,
			ExpectedRFC9457: &rfc9457.RFC9457{
				Type:   "about:blank",
				Title:  "failed to run query",
				Status: http.StatusInternalServerError,
			},
		},
		{
			Name: "missing type defaults to about:blank",
			JSON: `{"title":"failed to run query"}`,
			ExpectedRFC9457: &rfc9457.RFC9457{
				Type:  "about:blank",
				Title: "failed to run query",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			problem, err := rfc9457.FromJSON([]byte(tc.JSON))
			require.NoError(t, err)

			assert.Equal(t, tc.ExpectedRFC9457, problem)
		})
	}
}

func TestServeHTTPWritesProblem(t *testing.T) {
	problem := rfc9457.NewRFC9457(
		rfc9457.WithStatus(http.StatusTeapot),
		rfc9457.WithTitle("teapot"),
	)

	rec := httptest.NewRecorder()
	problem.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Equal(t, http.StatusTeapot, rec.Code)
	assert.Equal(t, rfc9457.ProblemContentType, rec.Header().Get("Content-Type"))

	body, err := io.ReadAll(rec.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"about:blank","title":"teapot","status":418}`, string(body))
}

func TestServeHTTPDefaultsStatusTo500(t *testing.T) {
	problem := rfc9457.NewRFC9457(rfc9457.WithTitle("oops"))

	rec := httptest.NewRecorder()
	problem.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, rfc9457.ProblemContentType, rec.Header().Get("Content-Type"))
}

func TestServeHTTPFallbackOnSerializationFailure(t *testing.T) {
	problem := rfc9457.NewRFC9457(
		rfc9457.WithStatus(http.StatusBadRequest),
		rfc9457.WithExtensions(rfc9457.NewExtension("title", "collides")),
	)

	rec := httptest.NewRecorder()
	problem.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, rfc9457.ProblemContentType, rec.Header().Get("Content-Type"))

	body, err := io.ReadAll(rec.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "failed to serialize problem details")
}

func TestFromError(t *testing.T) {
	problem := rfc9457.FromError(errors.New("boom"), http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, problem.Status)
	assert.Equal(t, "Bad Request", problem.Title)
	assert.Equal(t, "boom", problem.Detail)
}

func TestRFC9457ImplementsError(t *testing.T) {
	var err error = rfc9457.NewRFC9457(
		rfc9457.WithTitle("nope"),
		rfc9457.WithDetail("missing field"),
	)

	assert.EqualError(t, err, "rfc9457: nope: missing field")
}

func TestStatusNamedConstructors(t *testing.T) {
	cases := []struct {
		name   string
		got    *rfc9457.RFC9457
		status int
		title  string
	}{
		{"BadRequest", rfc9457.BadRequest(), http.StatusBadRequest, "Bad Request"},
		{"NotFound", rfc9457.NotFound(), http.StatusNotFound, "Not Found"},
		{"InternalServerError", rfc9457.InternalServerError(), http.StatusInternalServerError, "Internal Server Error"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.status, tc.got.Status)
			assert.Equal(t, tc.title, tc.got.Title)
		})
	}
}

func TestRoundTripPreservesExtensions(t *testing.T) {
	original := rfc9457.NewRFC9457(
		rfc9457.WithType("https://example.com/problems/oops"),
		rfc9457.WithStatus(http.StatusBadRequest),
		rfc9457.WithTitle("Validation failed"),
		rfc9457.WithDetail("missing email"),
		rfc9457.WithInstance("/users"),
		rfc9457.WithExtensions(
			rfc9457.NewExtension("trace_id", "abc-123"),
			rfc9457.NewExtension("retryable", true),
		),
	)

	body, err := original.ToJSON()
	require.NoError(t, err)

	parsed, err := rfc9457.FromJSON(body)
	require.NoError(t, err)

	assert.Equal(t, original.Type, parsed.Type)
	assert.Equal(t, original.Status, parsed.Status)
	assert.Equal(t, original.Title, parsed.Title)
	assert.Equal(t, original.Detail, parsed.Detail)
	assert.Equal(t, original.Instance, parsed.Instance)
	assert.Equal(t, "abc-123", parsed.Extensions["trace_id"])
	assert.Equal(t, true, parsed.Extensions["retryable"])
}

func TestFromJSONMalformed(t *testing.T) {
	_, err := rfc9457.FromJSON([]byte(`{not json`))
	require.Error(t, err)
	assert.ErrorIs(t, err, rfc9457.ErrUnableToUnmarshalJSON)
}

func TestFromErrorWithNilError(t *testing.T) {
	problem := rfc9457.FromError(nil, http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, problem.Status)
	assert.Equal(t, "Bad Request", problem.Title)
	assert.Equal(t, "", problem.Detail)
}

func TestErrorOnNilReceiver(t *testing.T) {
	var problem *rfc9457.RFC9457
	assert.Equal(t, "rfc9457: <nil>", problem.Error())
}

func TestStatusNamedConstructorAcceptsOverrides(t *testing.T) {
	problem := rfc9457.NotFound(
		rfc9457.WithDetail("user 42 not found"),
		rfc9457.WithInstance("/users/42"),
	)

	assert.Equal(t, http.StatusNotFound, problem.Status)
	assert.Equal(t, "Not Found", problem.Title)
	assert.Equal(t, "user 42 not found", problem.Detail)
	assert.Equal(t, "/users/42", problem.Instance)
}
