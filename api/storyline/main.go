package main

import (
	"fmt"
	commonmiddleware "server/api/middleware"
	storylinegin "server/api/storyline/ginhandler"
	storylineclient "server/api/storyline/serviceclient"
	"server/config"
	"time"

	"github.com/Jinchenyuan/wego"
	"github.com/Jinchenyuan/wego/logger"
	"github.com/Jinchenyuan/wego/transport"
	"github.com/Jinchenyuan/wego/transport/micro"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cfg, err := config.Read("config.toml")
	if err != nil {
		fmt.Printf("failed to read config: %v\n", err)
		return
	}

	m := wego.New(
		wego.WithEtcdConfig(clientv3.Config{
			Endpoints:   cfg.Etcd.Endpoints,
			DialTimeout: 5 * time.Second,
			Username:    cfg.Etcd.User,
			Password:    cfg.Etcd.Password,
		}),
		wego.WithHttpPort(cfg.Http.Port),
		wego.WithDSN(cfg.PostgreSQL.DSN),
		wego.WithLogLevel(logger.ParseLevel(cfg.Log.Level)),
		wego.WithRedisConfig(wego.RedisConfig{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
		wego.WithProfile(wego.Profile{
			Name: cfg.Profile.Name,
		}),
	)

	storylinegin.SetAuthMiddleware(commonmiddleware.AuthMiddleware(cfg.Http.ExcludeAuthPaths...))

	if err := storylinegin.Registry(); err != nil {
		fmt.Printf("failed to register storyline routes: %v\n", err)
		return
	}

	m.GetServerByType(transport.MICRO_SERVER).(*micro.Service).NewServiceClients(storylineclient.Registry)

	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}
