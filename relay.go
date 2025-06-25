package relay

import (
	"reflect"
)

type relay struct {
	handlers          map[any]AnyHandler
	middlewareHandler MiddlewareHandler
	defaultHandler    DefaultHandler

	messageHandlers          map[any]AnyMessageHandler
	middlewareMessageHandler MiddlewareMessageHandler
	defaultMessageHandler    DefaultMessageHandler
}

type Relay *relay

func requestKey[Request Req[Response], Response any]() any {
	return reflect.TypeFor[Request]()
}

func anyRequestKey(req any) any {
	return reflect.TypeOf(req)
}

func HandleAny(r *relay, req any) (any, error) {
	reqType := reflect.TypeOf(req)
	method, _ := reqType.MethodByName("ResponseType")
	resType := method.Type.Out(0)
	res := reflect.Zero(resType).Interface()
	var resErr error = nil
	key := anyRequestKey(req)
	handler, ok := r.handlers[key]
	anyCtx := NewAnyContext(
		func(a any) error {
			aType := reflect.TypeOf(a)
			if aType != reqType {
				return ErrInvalidType
			}
			req = a
			return nil
		},
		func() any { return req },
		func(a any) error {
			aType := reflect.TypeOf(a)
			if aType != resType {
				return ErrInvalidType
			}
			res = a
			return nil
		},
		func() any { return res },
		func(err error) { resErr = err },
		func() error { return resErr },
	)
	r.middlewareHandler(anyCtx, func() {
		if !ok {
			r.defaultHandler(anyCtx)
			return
		}
		res, err := handler(anyCtx.Req())
		anyCtx.SetRes(res)
		anyCtx.SetErr(err)
	})
	return res, resErr
}

func Handle[Request Req[Response], Response any](
	r *relay,
	req Request,
) (Response, error) {
	ctx := NewContext(req)
	key := requestKey[Request]()
	handler, ok := r.handlers[key]
	anyCtx := ctx.Any()
	r.middlewareHandler(anyCtx, func() {
		if !ok {
			r.defaultHandler(anyCtx)
			return
		}
		res, err := handler(ctx.Req())
		ctx.SetRes(res.(Response))
		ctx.SetErr(err)
	})
	return ctx.Res(), ctx.Err()
}

func messageKey[Mess Message]() any {
	return reflect.TypeFor[Mess]()
}

func anyMessageKey(message any) any {
	return reflect.TypeOf(message)
}

func HandleMessage[Mess Message](r *relay, m Mess) error {
	ctx := NewMessageCtx(m)
	key := messageKey[Mess]()
	handler, ok := r.messageHandlers[key]
	anyCtx := ctx.Any()
	r.middlewareMessageHandler(anyCtx, func() {
		if !ok {
			r.defaultMessageHandler(anyCtx)
			return
		}
		handler(ctx.Message())
	})
	return nil
}

func HandleAnyMessage(r *relay, message any) error {
	messType := reflect.TypeOf(message)
	key := anyMessageKey(message)
	handler, ok := r.messageHandlers[key]
	anyCtx := NewAnyMessageCtx(
		func(a Message) error {
			aType := reflect.TypeOf(a)
			if aType != messType {
				return ErrInvalidType
			}
			message = a
			return nil
		},
		func() Message { return message },
	)
	r.middlewareMessageHandler(anyCtx, func() {
		if !ok {
			r.defaultMessageHandler(anyCtx)
			return
		}
		handler(anyCtx.Message())
	})
	return nil
}
