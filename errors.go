package rfc9457

type RFC9457Error int

const (
	ErrUnableToMarshalJSON RFC9457Error = iota
	ErrUnableToTranslateToIntermediateMap
)

func (e RFC9457Error) Error() string {
	switch e {
	case ErrUnableToMarshalJSON:
		return "unable to marshal json"
	case ErrUnableToTranslateToIntermediateMap:
		return "unable to translate to intermediate map"
	default:
		return "unknown error"
	}
}
