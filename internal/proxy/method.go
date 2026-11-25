package proxy

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PATCH   Method = "PATCH"
	DELETE  Method = "DELETE"
	PUT     Method = "PUT"
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
	case string(DELETE):
		return DELETE, nil
	case string(PUT):
		return PUT, nil
	default:
		return INVALID, ErrInvalidMethod
	}
}
