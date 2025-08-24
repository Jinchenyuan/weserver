package main

import (
	"fmt"
	"server/core"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	m := core.New(
		core.WithEtcdConfig(clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5 * 1000,
			Username:    "root",
			Password:    "123456",
		}),
		core.WithHttpPort(8082),
		core.WithDSN("postgres://user:password@localhost:5432/land_contract?sslmode=disable"),
	)
	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}
