package relay

import (
	"reflect"
)

type relay struct {
	handlers       map[any]any
	defaultHandler DefaultHandler
}

type Relay *relay

func requestKey[Request Req[Response], Response any]() any {
	return reflect.TypeFor[Request]()
}

func Handle[Request Req[Response], Response any](
	r *relay,
	request Request,
) (Response, error) {
	key := requestKey[Request]()
	rawHandler, ok := r.handlers[key]
	if !ok {
		rawRes, err := r.defaultHandler(request)
		res, _ := rawRes.(Response)
		return res, err
	}
	handler := rawHandler.(Handler[Request, Response])
	return handler(request)
}
