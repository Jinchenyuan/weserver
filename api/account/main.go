package main

import (
	"fmt"
	"server/api/account/ginhandler"
	"server/api/account/serviceclient"
	"server/core"
	"server/core/transport"
	"server/core/transport/micro"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	m := core.New(
		core.WithEtcdConfig(clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5 * time.Second,
			Username:    "root",
			Password:    "123456",
		}),
		core.WithHttpPort(8083),
		core.WithDSN("postgres://user:password@localhost:5432/land_contract?sslmode=disable"),
	)
	core.SetGlobalMesa(m)

	ginhandler.Registry()

	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	ms.NewServiceClients(serviceclient.Registry)

	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}
