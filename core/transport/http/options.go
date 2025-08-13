package http

import (
	"net"
)

type Options func(o *options)

type options struct {
	Host net.IP
	Port int
}

func WithHost(host net.IP) Options {
	return func(o *options) {
		o.Host = host
	}
}

func WithPort(port int) Options {
	return func(o *options) {
		o.Port = port
	}
}
