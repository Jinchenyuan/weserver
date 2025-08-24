package core

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"server/account/model"
	"server/core/transport"
	"server/core/transport/http"
	"server/core/transport/micro"
	"server/third_part/etcd"
	"sync"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go-micro.dev/v5/registry"
	etcdReg "go-micro.dev/v5/registry/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Mesa struct {
	opts    options
	retChan chan int
	etcdCtl *etcd.Ctl
}

func New(opts ...Options) *Mesa {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	etcdCtl, err := etcd.NewCtl(etcd.ClientConfig{
		ConnectionType: etcd.ClientNonTLS,
		CertAuthority:  false,
		AutoTLS:        false,
		RevokeCerts:    false,
	}, etcd.WithEndpoints(o.EtcdConfig.Endpoints), etcd.WithAuth(o.EtcdConfig.Username, o.EtcdConfig.Password))
	if err != nil {
		panic(fmt.Sprintf("failed to create etcd controller: %v", err))
	}

	hs := http.NewHTTPServer(
		http.WithHost(net.ParseIP("0.0.0.0")),
		http.WithPort(o.HttpPort),
	)

	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   o.EtcdConfig.Endpoints,
		Username:    o.EtcdConfig.Username,
		Password:    o.EtcdConfig.Password,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create etcd clientv3: %v", err))
	}
	// Define a custom type for the context key to avoid collisions
	type etcdClientKeyType struct{}
	var etcdClientKey = etcdClientKeyType{}

	reg := etcdReg.NewEtcdRegistry(func(opt *registry.Options) {
		opt.Addrs = o.EtcdConfig.Endpoints
		opt.Context = context.WithValue(context.Background(), etcdClientKey, etcdCli)
	})
	ms := micro.NewMicroServer(micro.WithRegistry(reg))
	WithServers(hs, ms)(&o)

	return &Mesa{
		opts:    o,
		retChan: make(chan int),
		etcdCtl: etcdCtl,
	}
}

func (m *Mesa) Run() error {
	fmt.Println("start mesa!")

	m.startServers()

	m.addAccount()

	<-m.retChan

	return nil
}

func (m *Mesa) Stop() error {
	fmt.Println("stop mesa!")

	m.retChan <- 1

	for _, server := range m.opts.Servers {
		server.Stop(context.TODO())
	}

	m.etcdCtl.Close()

	return nil
}

func (m *Mesa) GetEtcdCtl() *etcd.Ctl {
	return m.etcdCtl
}

func (m *Mesa) startServers() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, server := range m.opts.Servers {
		wg.Add(1)
		go func(s transport.Server) {
			defer wg.Done()
			if err := s.Start(ctx); err != nil {
				fmt.Printf("server failed to start: %v\n", err)
				cancel() // Cancel context if any server fails
			}
		}(server)
	}

	wg.Wait()
	fmt.Println("All servers have started.")
}

func (m *Mesa) addAccount() {
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
