package rfc9457_test

import (
	"net/http"
	"testing"

	"alpineworks.io/rfc9457"
	"github.com/stretchr/testify/assert"
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
			ExpectedJSON: `{"type":"about:blank","title":"failed to run query","instance":""}`,
		},
		{
			Name: "with extensions",
			Options: []rfc9457.RFC9457Option{
				rfc9457.WithTitle("failed to run query"),
				rfc9457.WithExtensions(rfc9457.NewExtension("key", "value")),
			},
			ExpectedJSON: `{"type":"about:blank","title":"failed to run query","instance":"","key":"value"}`,
		},
		{
			Name: "with status",
			Options: []rfc9457.RFC9457Option{
				rfc9457.WithTitle("failed to run query"),
				rfc9457.WithStatus(http.StatusInternalServerError),
			},
			ExpectedJSON: `{"type":"about:blank","title":"failed to run query","instance":"","status":500}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			problem := rfc9457.NewRFC9457(tc.Options...)

			json, err := problem.ToJSON()
			assert.NoError(t, err)

			assert.JSONEq(t, tc.ExpectedJSON, json)
		})
	}
}

func TestJSONToRFC9457(t *testing.T) {
	tests := []struct {
		Name            string
		JSON            string
		ExpectedRFC9457 *rfc9457.RFC9457
	}{
		{
			Name: "simple",
			JSON: `{"type":"about:blank","title":"failed to run query","instance":""}`,
			ExpectedRFC9457: &rfc9457.RFC9457{
				Type:     "about:blank",
				Title:    "failed to run query",
				Instance: "",
			},
		},
		{
			Name: "with extensions",
			JSON: `{"type":"about:blank","title":"failed to run query","instance":"","key":"value"}`,
			ExpectedRFC9457: &rfc9457.RFC9457{
				Type:     "about:blank",
				Title:    "failed to run query",
				Instance: "",
				Extensions: map[string]any{
					"key": "value",
				},
			},
		},
		{
			Name: "with status",
			JSON: `{"type":"about:blank","title":"failed to run query","instance":"","status":500}`,
			ExpectedRFC9457: &rfc9457.RFC9457{
				Type:     "about:blank",
				Title:    "failed to run query",
				Instance: "",
				Status:   http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			problem, err := rfc9457.FromJSON(tc.JSON)
			assert.NoError(t, err)

			assert.Equal(t, tc.ExpectedRFC9457, problem)
		})
	}
}
