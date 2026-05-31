package proxy

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	INVALID Method = "INVALID"
)

func NewMethod(m string) (Method, error) {
	switch m {
	case string(GET):
		return GET, nil
	case string(POST):
		return POST, nil
	default:
		return INVALID, ErrInvalidMethod
	}
}
