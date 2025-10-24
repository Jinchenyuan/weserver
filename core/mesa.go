package core

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"server/core/transport"
	"server/core/transport/http"
	"server/core/transport/micro"
	"server/third_party/etcd"
	"sync"
	"syscall"

	"go-micro.dev/v5/registry"
	etcdReg "go-micro.dev/v5/registry/etcd"
)

type Mesa struct {
	opts          options
	retChan       chan int
	etcdCtl       *etcd.Ctl
	serversCtx    context.Context
	serversCancel context.CancelFunc
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
		http.WithType(transport.HTTP),
	)

	reg := etcdReg.NewEtcdRegistry(
		registry.Addrs(o.EtcdConfig.Endpoints...),
		etcdReg.Auth(o.EtcdConfig.Username, o.EtcdConfig.Password),
	)
	ms := micro.NewMicroServer(
		micro.WithRegistry(reg),
		micro.WithType(transport.MICRO_SERVER),
		micro.WithServiceScheme(o.serviceScheme),
	)
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

	go m.waitForStop()

	<-m.retChan

	return nil
}

func (m *Mesa) GetServerByType(typ transport.NetType) transport.Server {
	for _, server := range m.opts.Servers {
		s := server.GetType()
		if s == typ {
			return server
		}
	}
	return nil
}

func (m *Mesa) startServers() {
	var wg sync.WaitGroup
	m.serversCtx, m.serversCancel = context.WithCancel(context.Background())

	for _, server := range m.opts.Servers {
		wg.Add(1)
		go func(s transport.Server) {
			defer wg.Done()
			if err := s.Start(m.serversCtx); err != nil {
				fmt.Printf("server failed to start: %v\n", err)
				m.serversCancel() // Cancel context if any server fails
			}
		}(server)
	}

	wg.Wait()
	fmt.Println("All servers have started.")
}

func (m *Mesa) waitForStop() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	sig := <-sigChan
	fmt.Printf("Received signal: %v. Shutting down servers...\n", sig)

	m.serversCancel()

	m.etcdCtl.Close()

	m.retChan <- 1

	fmt.Println("Mesa has shut down.")
}
