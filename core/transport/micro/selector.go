package micro

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"go-micro.dev/v5/metadata"
	"go-micro.dev/v5/registry"
	"go-micro.dev/v5/registry/cache"
	"go-micro.dev/v5/selector"
)

type idSelector struct {
	so selector.Options
	rc cache.Cache
	mu sync.RWMutex
}

func (c *idSelector) newCache() cache.Cache {
	opts := make([]cache.Option, 0, 1)

	if c.so.Context != nil {
		if t, ok := c.so.Context.Value("selector_ttl").(time.Duration); ok {
			opts = append(opts, cache.WithTTL(t))
		}
	}

	return cache.New(c.so.Registry, opts...)
}

func (c *idSelector) Init(opts ...selector.Option) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, o := range opts {
		o(&c.so)
	}

	c.rc.Stop()
	c.rc = c.newCache()

	return nil
}

func (c *idSelector) Options() selector.Options {
	return c.so
}

func (c *idSelector) Select(service string, opts ...selector.SelectOption) (selector.Next, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	sopts := selector.SelectOptions{
		Strategy: c.so.Strategy,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	// get the service
	// try the cache first
	// if that fails go directly to the registry
	services, err := c.rc.GetService(service)
	if err != nil {
		if errors.Is(err, registry.ErrNotFound) {
			return nil, selector.ErrNotFound
		}

		return nil, err
	}

	if sopts.Context != nil {
		if md, ok := metadata.FromContext(sopts.Context); ok {
			// 打印md
			fmt.Printf("metadata: %v\n", md)
		}
	}

	// apply the filters
	for _, filter := range sopts.Filters {
		services = filter(services)
	}

	// if there's nothing left, return
	if len(services) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	// 将 services 展平成 nodes 列表（示例）
	var nodes []*registry.Node
	for _, svc := range services {
		nodes = append(nodes, svc.Nodes...)
	}

	for _, n := range nodes {
		// print node id
		fmt.Printf("node id: %s\n", n.Id)
	}

	// TODO: 直接实现一个闭包策略，不允许外部传入策略
	return sopts.Strategy(services), nil
}

func (c *idSelector) Mark(service string, node *registry.Node, err error) {
}

func (c *idSelector) Reset(service string) {
}

// Close stops the watcher and destroys the cache.
func (c *idSelector) Close() error {
	c.rc.Stop()

	return nil
}

func (c *idSelector) String() string {
	return "idRegistry"
}

func newSelector(opts ...selector.Option) selector.Selector {
	sopts := selector.Options{
		Strategy: selector.Random,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	if sopts.Registry == nil {
		sopts.Registry = registry.DefaultRegistry
	}

	s := &idSelector{
		so: sopts,
	}
	s.rc = s.newCache()

	return s
}

func NewSelectorDependId(reg registry.Registry) selector.Selector {
	return newSelector(
		selector.SetStrategy(selector.Random),
		selector.Registry(reg),
	)
}
