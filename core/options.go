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

	// AppId are provided for logical separation, but clusters across AppId are isolated with no built-in communication.
	// Users requiring interaction must handle it externally.
	AppId string

	EtcdConfig clientv3.Config

	Servers []transport.Server
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

func WithAppId(appId string) Options {
	return func(o *options) {
		o.AppId = appId
	}
}

func WithEtcdConfig(config clientv3.Config) Options {
	return func(o *options) {
		o.EtcdConfig = config
	}
}
