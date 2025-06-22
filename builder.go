package relay

import (
	"errors"
	"fmt"
)

type Builder struct {
	valid       bool
	middlewares *[]MiddlewareHandler
	r           *relay
}

var (
	ErrDidntUseCtor         error = errors.New("use constructor")
	ErrHandlerAlreadyExists error = errors.New("handler already exists")
	ErrHandlerNotFound      error = errors.New("handler wasn't found")
)

func NewBuilder() Builder {
	return Builder{
		valid:       true,
		middlewares: &[]MiddlewareHandler{},
		r: &relay{
			handlers: map[any]any{},
			defaultHandler: func(req AnyContext) {
				req.SetErr(ErrHandlerNotFound)
			},
		},
	}
}

func (b Builder) Wrap(wrapped func(Builder) Builder) Builder {
	if !b.valid {
		panic(fmt.Sprintf("%s\n", ErrDidntUseCtor.Error()))
	}
	return wrapped(b)
}

func (b Builder) DefaultHandler(handler DefaultHandler) Builder {
	if !b.valid {
		panic(fmt.Sprintf("%s\n", ErrDidntUseCtor.Error()))
	}
	b.r.defaultHandler = handler
	return b
}

func (b Builder) RegisterMiddleware(handler MiddlewareHandler) Builder {
	(*b.middlewares) = append((*b.middlewares), handler)
	return b
}

// can panic. if you do not want to panic use TryRegister
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
	b.r.handlers[key] = handler
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
	r := *b.r
	return &r
}
