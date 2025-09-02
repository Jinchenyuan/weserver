package main

import (
	"fmt"
	"server/api/account/ginhandler"
	"server/api/account/serviceclient"
	"server/core"
	"server/core/config"
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
		core.WithHttpPort(cfg.HTTP.Port),
		core.WithDSN(cfg.PostgreSQL.DSN),
	)
	core.SetGlobalMesa(m)

	ginhandler.Registry()

	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	ms.NewServiceClients(serviceclient.Registry)

	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}
