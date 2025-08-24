package micro

import "go-micro.dev/v5/registry"

type Options func(o *options)

type options struct {
	reg registry.Registry
}

func WithRegistry(reg registry.Registry) Options {
	return func(o *options) {
		o.reg = reg
	}
}
