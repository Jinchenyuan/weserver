package micro

import (
	"context"
	"fmt"
	"log"
	"server/core/transport"

	goMicro "go-micro.dev/v5"
	"go-micro.dev/v5/registry"
)

type RegisterHandler func(goMicro.Service) error
type NewServiceClients func(reg registry.Registry) map[string]any

type Service struct {
	opts    options
	clients map[string]any
}

func NewMicroServer(opts ...Options) *Service {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	ms := &Service{
		opts: o,
	}

	return ms
}

func (s *Service) GetServiceClient(service transport.ServiceType) any {
	if s.clients == nil {
		return nil
	}
	if _, ok := s.clients[string(service)]; !ok {
		return nil
	}
	return s.clients[string(service)]
}

func (s *Service) NewServiceClients(nsc NewServiceClients) {
	s.clients = nsc(s.opts.reg)
}

func (s *Service) RegisterHandler(handler RegisterHandler) {
	goService := goMicro.NewService(
		goMicro.Name(string(s.opts.serviceScheme.Name)),
		goMicro.Address(fmt.Sprintf(":%d", s.opts.serviceScheme.Port)),
		goMicro.Registry(s.opts.reg),
	)
	goService.Init()
	if err := handler(goService); err != nil {
		log.Fatalf("RegisterHandler err: %v", err)
	}

	go func() {
		if err := goService.Run(); err != nil {
			log.Fatalf("goService.Run err: %v", err)
		}
	}()

}

func (s *Service) GetType() transport.NetType {
	return s.opts.Type
}

func (s *Service) Start(ctx context.Context) error {
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	return nil
}
