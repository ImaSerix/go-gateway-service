package endpoint

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PATCH   Method = "PATCH"
	INVALID Method = "INVALID"
)

func NewMethod(m string) (Method, error) {
	switch m {
	case string(GET):
		return GET, nil
	case string(POST):
		return POST, nil
	case string(PATCH):
		return PATCH, nil
	default:
		return INVALID, ErrInvalidMethod
	}
}
