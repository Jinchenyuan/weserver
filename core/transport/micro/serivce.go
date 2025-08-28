package micro

import (
	"context"
	"fmt"
	"log"
	"server/core/transport"

	goMicro "go-micro.dev/v5"
)

type RegisterHandler func(goMicro.Service) error

type Service struct {
	opts options
}

func NewMicroServer(opts ...Options) *Service {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	ms := &Service{
		opts: options{},
	}

	return ms
}

func (s *Service) RegisterHandler(handler RegisterHandler) {
	goService := goMicro.NewService(
		goMicro.Name(s.opts.serviceScheme.Name),
		goMicro.Address(fmt.Sprintf(":%d", s.opts.serviceScheme.Port)),
		goMicro.Registry(s.opts.reg),
	)
	goService.Init()
	if err := handler(goService); err != nil {
		log.Fatalf("RegisterHandler err: %v", err)
	}

	if err := goService.Run(); err != nil {
		log.Fatalf("goService.Run err: %v", err)
	}
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
