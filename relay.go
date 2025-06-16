package relay

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type relay struct {
	registerMutex sync.Mutex
	handlers      map[any]any
	config        Config
}

type Relay *relay

// can return errors: ErrInvalidConfig
func TryNewRelay(config Config) (Relay, error) {
	if !config.valid {
		return nil, ErrInvalidConfig
	}
	return &relay{
		registerMutex: sync.Mutex{},
		handlers:      make(map[any]any),
		config:        config,
	}, nil
}

func NewRelay(config Config) Relay {
	relay, err := TryNewRelay(config)
	if err != nil {
		panic(fmt.Sprintf("%s\n", err))
	}
	return relay
}

func requestKey[Request Req[Response], Response any]() any {
	return reflect.TypeFor[Request]()
}

var (
	ErrHandlerAlreadyExists error = errors.New("handler already exists")
	ErrHandlerNotFound      error = errors.New("handler wasn't found")
)

// can return ErrHandlerAlreadyExists
func TryRegister[Request Req[Response], Response any](
	r *relay,
	handler Handler[Request, Response],
) error {
	key := requestKey[Request]()
	r.registerMutex.Lock()
	defer r.registerMutex.Unlock()

	// verify do handler exists
	if _, ok := r.handlers[key]; ok {
		return ErrHandlerAlreadyExists
	}

	// reigster
	r.handlers[key] = handler
	return nil
}

// can panic. if you do not want to panic use TryRegister
func Register[Request Req[Response], Response any](
	r *relay,
	handler Handler[Request, Response],
) {
	if err := TryRegister(r, handler); err != nil {
		panic(fmt.Sprintf("%s\n", err.Error()))
	}
}

func Handle[Request Req[Response], Response any](
	r *relay,
	request Request,
) (Response, error) {
	key := requestKey[Request]()
	rawHandler, ok := r.handlers[key]
	if !ok {
		rawRes, err := r.config.defaultHandler(request)
		res, _ := rawRes.(Response)
		return res, err
	}
	handler := rawHandler.(Handler[Request, Response])
	return handler(request)
}
