package relay

import (
	"reflect"
)

type relay struct {
	handlers          map[any]any
	middlewareHandler MiddlewareHandler
	defaultHandler    DefaultHandler
}

type Relay *relay

func requestKey[Request Req[Response], Response any]() any {
	return reflect.TypeFor[Request]()
}

func Handle[Request Req[Response], Response any](
	r *relay,
	req Request,
) (Response, error) {
	ctx := NewContext(req)
	key := requestKey[Request]()
	rawHandler, ok := r.handlers[key]
	anyCtx := ctx.Any()
	r.middlewareHandler(anyCtx, func() {
		if !ok {
			r.defaultHandler(anyCtx)
			return
		}
		handler := rawHandler.(Handler[Request, Response])
		res, err := handler(ctx.Req())
		ctx.SetRes(res)
		ctx.SetErr(err)
	})
	return ctx.Res(), ctx.Err()
}
