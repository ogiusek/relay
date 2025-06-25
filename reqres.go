package relay

// e.g.
//
//	type EgRes struct {
//		ExampleField int
//	}
type Res interface{}

// every interface should use this
// e.g.
//
//	type EgReq struct {
//		Req[EgRes] // it can be nil its just a marker
//		Field1     int
//		Field2     int
//	}
type Req[Response Res] interface {
	// this is never called
	// this is just request marker
	ResponseType() Response
}

// everything below should be treated as if it had prefix "Request"

type AnyHandler func(any) (any, error)
type Handler[Request Req[Response], Response any] func(Request) (Response, error)
type DefaultHandler func(ctx AnyContext)
type MiddlewareHandler func(ctx AnyContext, next func())

//

type Context[Request Req[Response], Response Res] interface {
	SetReq(Request)
	Req() Request

	SetRes(Response)
	Res() Response

	SetErr(error)
	Err() error
	Any() AnyContext
}

type context[Request Req[Response], Response Res] struct {
	req Request
	res Response
	err error
}

func NewContext[Request Req[Response], Response Res](req Request) Context[Request, Response] {
	return &context[Request, Response]{req: req}
}

func (ctx *context[Request, Response]) SetReq(req Request)  { ctx.req = req }
func (ctx *context[Request, Response]) Req() Request        { return ctx.req }
func (ctx *context[Request, Response]) SetRes(res Response) { ctx.res = res }
func (ctx *context[Request, Response]) Res() Response       { return ctx.res }
func (ctx *context[Request, Response]) SetErr(err error)    { ctx.err = err }
func (ctx *context[Request, Response]) Err() error          { return ctx.err }
func (ctx *context[Request, Response]) Any() AnyContext {
	return NewAnyContext(
		func(rawReq any) error {
			req, ok := rawReq.(Request)
			if !ok {
				return ErrInvalidType
			}
			ctx.req = req
			return nil
		},
		func() any { return ctx.req },
		func(rawRes any) error {
			res, ok := rawRes.(Response)
			if !ok {
				return ErrInvalidType
			}
			ctx.res = res
			return nil
		},
		func() any { return ctx.res },
		ctx.SetErr,
		ctx.Err,
	)
}

//

type AnyContext interface {
	Req() any
	SetReq(any) error

	// can return type errors
	SetRes(any) error
	Res() any

	SetErr(error)
	Err() error
}

type anyContext struct {
	setReq func(any) error
	req    func() any
	setRes func(any) error
	res    func() any
	setErr func(error)
	err    func() error
}

func NewAnyContext(
	setReq func(any) error,
	req func() any,
	setRes func(any) error,
	res func() any,
	setErr func(error),
	err func() error,
) AnyContext {
	return anyContext{
		setReq: setReq,
		req:    req,
		setRes: setRes,
		res:    res,
		setErr: setErr,
		err:    err,
	}

}

func (ctx anyContext) SetReq(req any) error { return ctx.setReq(req) }
func (ctx anyContext) Req() any             { return ctx.req() }
func (ctx anyContext) SetRes(res any) error { return ctx.setRes(res) }
func (ctx anyContext) Res() any             { return ctx.res() }
func (ctx anyContext) SetErr(err error)     { ctx.setErr(err) }
func (ctx anyContext) Err() error           { return ctx.err() }
