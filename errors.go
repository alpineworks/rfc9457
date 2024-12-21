package rfc9457

type RFC9457Error int

const (
	ErrUnableToMarshalJSON RFC9457Error = iota
	ErrUnableToUnmarshalJSON
	ErrUnableToTranslateToIntermediateMap
	ErrUnableToTranslateToRFC9457
)

func (e RFC9457Error) Error() string {
	switch e {
	case ErrUnableToMarshalJSON:
		return "unable to marshal json"
	case ErrUnableToUnmarshalJSON:
		return "unable to unmarshal json"
	case ErrUnableToTranslateToIntermediateMap:
		return "unable to translate to intermediate map"
	case ErrUnableToTranslateToRFC9457:
		return "unable to translate to rfc9457"
	default:
		return "unknown error"
	}
}
