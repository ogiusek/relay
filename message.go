package relay

// message differs from request in that it does not have a response

type Message interface{}

type MessageHandler[Mess Message] func(Mess)
type DefaultMessageHandler func(ctx AnyMessageCtx)
type MiddlewareMessageHandler func(ctx AnyMessageCtx, next func())

type MessageCtx[Mess Message] interface {
	Message() Mess
	SetMessage(Mess)

	Any() AnyMessageCtx
}

type messageCtx[Mess Message] struct {
	message Mess
}

func NewMessageCtx[Mess Message](message Mess) MessageCtx[Mess] {
	return &messageCtx[Mess]{message: message}
}

func (mess *messageCtx[Mess]) Message() Mess           { return mess.message }
func (mess *messageCtx[Mess]) SetMessage(message Mess) { mess.message = message }
func (mess *messageCtx[Mess]) Any() AnyMessageCtx {
	return NewAnyMessageCtx(
		func() Message { return mess.message },
		func(m Message) error {
			message, ok := m.(Mess)
			if !ok {
				return ErrInvalidType
			}
			mess.message = message
			return nil
		},
	)
}

type AnyMessageCtx interface {
	Message() Message
	// ErrInvalidType
	SetMessage(Message) error
}

type anyMessageCtx struct {
	message    func() Message
	setMessage func(Message) error
}

func NewAnyMessageCtx(
	message func() Message,
	setMessage func(Message) error,
) AnyMessageCtx {
	return anyMessageCtx{
		message:    message,
		setMessage: setMessage,
	}
}

func (ctx anyMessageCtx) Message() Message              { return ctx.message() }
func (ctx anyMessageCtx) SetMessage(mess Message) error { return ctx.setMessage(mess) }
