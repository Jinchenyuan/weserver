package core

import (
	"server/core/logger"
	"server/core/transport"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Options func(o *options)

type options struct {
	// log level
	// default: common.Info
	// others: common.Debug, common.Warn, common.Error, common.Fatal.
	// see common.Level
	LogLevel logger.Level

	HttpPort int

	EtcdConfig clientv3.Config

	dsn string

	Servers []transport.Server
}

func WithDSN(dsn string) Options {
	return func(o *options) {
		o.dsn = dsn
	}
}

func WithHttpPort(port int) Options {
	return func(o *options) {
		o.HttpPort = port
	}
}

func WithLogLevel(level logger.Level) Options {
	return func(o *options) {
		o.LogLevel = level
	}
}

func WithServers(servers ...transport.Server) Options {
	return func(o *options) {
		o.Servers = servers
	}
}

func WithEtcdConfig(config clientv3.Config) Options {
	return func(o *options) {
		o.EtcdConfig = config
	}
}
