package main

import (
	"fmt"
	"server/config"
	"server/service/account/servicehandler"
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
		wego.WithServiceScheme(micro.ServiceScheme{
			Name:    cfg.Service.Name,
			Version: cfg.Service.Version,
			Port:    cfg.Service.Port,
		}),
		wego.WithProfile(wego.Profile{
			Name: cfg.Profile.Name,
		}),
	)

	microService := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	microService.RegisterHandler(servicehandler.Registry)

	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}
