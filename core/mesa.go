package core

import (
	"context"
	"fmt"
	"server/core/transport"
	"server/third_part/etcd"
	"sync"
)

type Mesa struct {
	opts        options
	retChan     chan int
	etcdCtl     *etcd.Ctl
	clusterId   int64
	clusterList []int64
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
		panic(fmt.Sprintf("failed to create etcd client: %v", err))
	}

	return &Mesa{
		opts:        o,
		retChan:     make(chan int),
		etcdCtl:     etcdCtl,
		clusterList: make([]int64, 0),
	}
}

func (m *Mesa) Run() error {
	fmt.Println("start mesa!")

	m.startServers()

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
