package micro

import (
	"context"
	"server/core/transport"
)

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

func (s *Service) GetType() transport.NetType {
	return s.opts.Type
}

func (s *Service) Start(ctx context.Context) error {
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	return nil
}
