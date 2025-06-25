package relay

import (
	"errors"
	"fmt"
)

type Builder struct {
	valid              bool
	middlewares        *[]MiddlewareHandler
	messageMiddlewares *[]MiddlewareMessageHandler
	r                  *relay
}

var (
	ErrDidntUseCtor         error = errors.New("use constructor")
	ErrHandlerAlreadyExists error = errors.New("handler already exists")
	ErrHandlerNotFound      error = errors.New("handler wasn't found")
)

func NewBuilder() Builder {
	return Builder{
		valid:              true,
		middlewares:        &[]MiddlewareHandler{},
		messageMiddlewares: &[]MiddlewareMessageHandler{},
		r: &relay{
			handlers:        map[any]AnyHandler{},
			messageHandlers: map[any]AnyMessageHandler{},
			defaultHandler: func(req AnyContext) {
				req.SetErr(ErrHandlerNotFound)
			},
			defaultMessageHandler: func(ctx AnyMessageCtx) {},
		},
	}
}

func (b Builder) Wrap(wrapped func(Builder) Builder) Builder {
	if !b.valid {
		panic(fmt.Sprintf("%s\n", ErrDidntUseCtor.Error()))
	}
	return wrapped(b)
}

func (b Builder) DefaultMessageHandler(handler DefaultMessageHandler) Builder {
	if !b.valid {
		panic(fmt.Sprintf("%s\n", ErrDidntUseCtor.Error()))
	}
	b.r.defaultMessageHandler = handler
	return b
}

func (b Builder) DefaultHandler(handler DefaultHandler) Builder {
	if !b.valid {
		panic(fmt.Sprintf("%s\n", ErrDidntUseCtor.Error()))
	}
	b.r.defaultHandler = handler
	return b
}

func (b Builder) RegisterMessageMiddleware(handler MiddlewareMessageHandler) Builder {
	(*b.messageMiddlewares) = append((*b.messageMiddlewares), handler)
	return b
}

func (b Builder) RegisterMiddleware(handler MiddlewareHandler) Builder {
	(*b.middlewares) = append((*b.middlewares), handler)
	return b
}

func MessageRegister[Mess Message](
	b Builder,
	handler MessageHandler[Mess],
) Builder {
	if !b.valid {
		panic(fmt.Sprintf("%s\n", ErrDidntUseCtor.Error()))
	}
	key := messageKey[Mess]()

	// verify do handler exists
	if _, ok := b.r.handlers[key]; ok {
		panic(fmt.Sprintf("%s\n", ErrHandlerAlreadyExists.Error()))
	}

	// reigster
	h := handler
	// b.r.messageHandlers[key] = handler
	b.r.messageHandlers[key] = func(m Message) error {
		mess, ok := m.(Mess)
		if !ok {
			return ErrInvalidType
		}
		h(mess)
		return nil
	}
	return b
}

func Register[Request Req[Response], Response any](
	b Builder,
	handler Handler[Request, Response],
) Builder {
	if !b.valid {
		panic(fmt.Sprintf("%s\n", ErrDidntUseCtor.Error()))
	}
	key := requestKey[Request]()

	// verify do handler exists
	if _, ok := b.r.handlers[key]; ok {
		panic(fmt.Sprintf("%s\n", ErrHandlerAlreadyExists.Error()))
	}

	// reigster
	h := handler
	b.r.handlers[key] = func(a any) (any, error) {
		req, ok := a.(Request)
		if !ok {
			var res Response
			return res, ErrInvalidType
		}
		return h(req)
	}
	return b
}

func (b Builder) Build() Relay {
	if !b.valid {
		panic(fmt.Sprintf("%s\n", ErrDidntUseCtor.Error()))
	}
	var m MiddlewareHandler = func(req AnyContext, next func()) { next() }
	for i, middleware := range *b.middlewares {
		if i == 0 {
			m = middleware
			continue
		}
		wrappedM := m
		m = func(req AnyContext, next func()) {
			middleware(req, func() { wrappedM(req, next) })
		}
	}
	b.r.middlewareHandler = m
	var mM MiddlewareMessageHandler = func(ctx AnyMessageCtx, next func()) { next() }
	for i, middleware := range *b.messageMiddlewares {
		if i == 0 {
			mM = middleware
			continue
		}
		wrappedM := mM
		mM = func(req AnyMessageCtx, next func()) {
			middleware(req, func() { wrappedM(req, next) })
		}
	}
	b.r.middlewareMessageHandler = mM
	r := *b.r
	return &r
}
