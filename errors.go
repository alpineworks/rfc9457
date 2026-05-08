package rfc9457

import "errors"

var (
	ErrUnableToMarshalJSON   = errors.New("rfc9457: unable to marshal json")
	ErrUnableToUnmarshalJSON = errors.New("rfc9457: unable to unmarshal json")
	ErrExtensionKeyCollision = errors.New("rfc9457: extension key collides with reserved field")
)
