package main

import (
	"context"
	"database/sql"
	"fmt"
	"server/core"
	"server/core/transport"
	"server/core/transport/micro"
	"server/model"
	"server/service/account/servicehandler"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
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
		core.WithHttpPort(8093),
		core.WithDSN("postgres://user:password@localhost:5432/land_contract?sslmode=disable"),
		core.WithServiceScheme(micro.ServiceScheme{
			Name:    transport.Account,
			Version: "v1",
			Port:    8193,
		}),
	)

	microService := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	microService.RegisterHandler(servicehandler.Registry)

	addAccount()

	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}

func addAccount() {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://user:password@localhost:5432/land_contract?sslmode=disable")))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	account := &model.Account{
		OwnerID:   1,
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	account.SetDB(db)
	err := account.Create(context.Background())
	if err != nil {
		fmt.Printf("failed to create account: %v\n", err)
	}
}
