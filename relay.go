package relay

import "reflect"

type relay struct {
	handlers          map[any]any
	middlewareHandler MiddlewareHandler
	defaultHandler    DefaultHandler

	messageHandlers          map[any]any
	middlewareMessageHandler MiddlewareMessageHandler
	defaultMessageHandler    DefaultMessageHandler
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

func messageKey[Mess Message]() any {
	return reflect.TypeFor[Mess]()
}

func HandleMessage[Mess Message](r *relay, m Mess) {
	ctx := NewMessageCtx(m)
	key := messageKey[Mess]()
	rawHandler, ok := r.messageHandlers[key]
	anyCtx := ctx.Any()
	r.middlewareMessageHandler(anyCtx, func() {
		if !ok {
			r.defaultMessageHandler(anyCtx)
			return
		}
		handler := rawHandler.(MessageHandler[Mess])
		handler(ctx.Message())
	})
}
