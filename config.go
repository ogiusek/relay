package relay

import "errors"

var (
	ErrInvalidConfig error = errors.New("invalid cofig. use constructor")
)

type Config struct {
	valid          bool
	defaultHandler DefaultHandler
}

func NewConfig(
	defaultHandler DefaultHandler,
) Config {
	return Config{
		valid:          true,
		defaultHandler: defaultHandler,
	}
}

type ConfigBuilder interface {
	DefaultHandler(handler DefaultHandler) ConfigBuilder
	Build() Config
}

type configBuilder struct {
	Config
}

func NewConfigBuilder() ConfigBuilder {
	return configBuilder{
		Config: NewConfig(func(a any) (Res, error) { return nil, ErrHandlerNotFound }),
	}
}

func (builder configBuilder) DefaultHandler(handler DefaultHandler) ConfigBuilder {
	builder.defaultHandler = handler
	return builder
}

func (builder configBuilder) Build() Config {
	return builder.Config
}
