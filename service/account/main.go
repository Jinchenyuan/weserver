package main

import (
	"context"
	"database/sql"
	"fmt"
	"server/core"
	"server/core/config"
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
		core.WithServiceScheme(micro.ServiceScheme{
			Name:    transport.ServiceType(cfg.Service.Name),
			Version: cfg.Service.Version,
			Port:    cfg.Service.Port,
		}),
	)

	core.SetGlobalMesa(m)

	microService := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	microService.RegisterHandler(servicehandler.Registry)

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
