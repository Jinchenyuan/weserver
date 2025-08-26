package micro

import (
	"server/core/transport"

	"go-micro.dev/v5/registry"
)

type Options func(o *options)

type options struct {
	reg  registry.Registry
	Type transport.NetType
}

func WithType(typ transport.NetType) Options {
	return func(o *options) {
		o.Type = typ
	}
}

func WithRegistry(reg registry.Registry) Options {
	return func(o *options) {
		o.reg = reg
	}
}
