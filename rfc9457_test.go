package rfc9457_test

import (
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
