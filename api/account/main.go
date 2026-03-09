package main

import (
	"context"
	"fmt"
	"server/api/account/ginhandler"
	"server/api/account/serviceclient"
	"server/core"
	"server/core/config"
	"server/core/logger"
	"server/core/middleware"
	"server/core/transport"
	"server/core/transport/micro"
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
		core.WithProfile(core.Profile{
			Name: cfg.Profile.Name,
		}),
	)

	ginhandler.SetAuthMiddleware(middleware.AuthMiddleware("account", func(id string) string {
		cacheToken, err := m.Redis.Get(context.Background(), fmt.Sprintf("token:%s", id)).Result()
		if err != nil {
			return ""
		}
		return cacheToken
	}, cfg.Http.ExcludeAuthPaths...))

	ginhandler.Registry()

	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	ms.NewServiceClients(serviceclient.Registry)

	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}
