package main

import (
	"fmt"
	"server/core"
	"server/core/config"
	"server/core/logger"
	"server/core/transport"
	"server/core/transport/micro"
	"server/service/account/servicehandler"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cfg, err := config.Read("config.toml")
	if err != nil {
		fmt.Printf("failed to read config: %v\n", err)
		return
	}

	m := core.New(
		core.WithEtcdConfig(clientv3.Config{
			Endpoints:   cfg.Etcd.Endpoints,
			DialTimeout: 5 * time.Second,
			Username:    cfg.Etcd.User,
			Password:    cfg.Etcd.Password,
		}),
		core.WithHttpPort(cfg.Http.Port),
		core.WithDSN(cfg.PostgreSQL.DSN),
		core.WithLogLevel(logger.ParseLevel(cfg.Log.Level)),
		core.WithRedisConfig(core.RedisConfig{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
		core.WithServiceScheme(micro.ServiceScheme{
			Name:    transport.ServiceType(cfg.Service.Name),
			Version: cfg.Service.Version,
			Port:    cfg.Service.Port,
		}),
		core.WithProfile(core.Profile{
			Name: cfg.Profile.Name,
		}),
	)

	microService := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	microService.RegisterHandler(servicehandler.Registry)

	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}
