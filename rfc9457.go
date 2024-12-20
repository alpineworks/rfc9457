package rfc9457

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

type RFC9457 struct {
	// The "type" member is a JSON string containing a URI reference [URI] that identifies the problem type.
	// Consumers MUST use the "type" URI (after resolution, if necessary) as the problem type's primary identifier.
	// When this member is not present, its value is assumed to be "about:blank".
	//
	// https://datatracker.ietf.org/doc/html/rfc9457#name-type
	Type string `json:"type" mapstructure:"type"`

	// The "status" member is a JSON number indicating the HTTP status code generated by the origin server for this occurrence of the problem.
	// The "status" member, if present, is only advisory; it conveys the HTTP status code used for the convenience of the consumer.
	//
	// https://datatracker.ietf.org/doc/html/rfc9457#name-status
	Status string `json:"status,omitempty" mapstructure:"status,omitempty"`

	// The "title" member is a JSON string containing a short, human-readable summary of the problem type.
	// It SHOULD NOT change from occurrence to occurrence of the problem, except for localization.
	// The "title" string is advisory and is included only for users who are unaware of and cannot discover the semantics of the type URI.
	//
	// https://datatracker.ietf.org/doc/html/rfc9457#name-title
	Title string `json:"title" mapstructure:"title"`

	// The "detail" member is a JSON string containing a human-readable explanation specific to this occurrence of the problem.
	// The "detail" string, if present, ought to focus on helping the client correct the problem, rather than giving debugging information.
	//
	// https://datatracker.ietf.org/doc/html/rfc9457#name-detail
	Detail string `json:"detail,omitempty" mapstructure:"detail,omitempty"`

	// The "instance" member is a JSON string containing a URI reference that identifies the specific occurrence of the problem.
	//
	// https://datatracker.ietf.org/doc/html/rfc9457#name-instance
	Instance string `json:"instance" mapstructure:"instance"`

	// Problem type definitions MAY extend the problem details object with additional members that are specific to that problem type.
	//
	// https://datatracker.ietf.org/doc/html/rfc9457#name-extension-members
	extensions map[string]any `json:"-"`
}

type RFC9457Option func(*RFC9457)

func NewRFC9457(options ...RFC9457Option) *RFC9457 {
	r := &RFC9457{
		Type: "about:blank",
	}

	// apply options
	for _, o := range options {
		o(r)
	}

	return r
}

func WithType(t string) RFC9457Option {
	return func(r *RFC9457) {
		r.Type = t
	}
}

func WithStatus(s string) RFC9457Option {
	return func(r *RFC9457) {
		r.Status = s
	}
}

func WithTitle(t string) RFC9457Option {
	return func(r *RFC9457) {
		r.Title = t
	}
}

func WithDetail(d string) RFC9457Option {
	return func(r *RFC9457) {
		r.Detail = d
	}
}

func WithInstance(i string) RFC9457Option {
	return func(r *RFC9457) {
		r.Instance = i
	}
}

type Extension struct {
	Key   string
	Value any
}

func NewExtension(key string, value any) Extension {
	return Extension{
		Key:   key,
		Value: value,
	}
}

func WithExtensions(e ...Extension) RFC9457Option {
	extensions := make(map[string]any)
	for _, ext := range e {
		extensions[ext.Key] = ext.Value
	}

	return func(r *RFC9457) {
		r.extensions = extensions
	}
}

func (r *RFC9457) ToJSON() (string, error) {
	intermediateMap := make(map[string]any)

	err := mapstructure.Decode(r, &intermediateMap)
	if err != nil {
		return "", ErrUnableToTranslateToIntermediateMap
	}

	// set extensions
	//
	// will not override existing keys (as they are part of the rfc spec)
	for k, v := range r.extensions {
		if _, exists := intermediateMap[k]; !exists {
			intermediateMap[k] = v
		}
	}

	jsonResult, err := json.Marshal(intermediateMap)
	if err != nil {
		return "", ErrUnableToMarshalJSON
	}

	return string(jsonResult), nil
}