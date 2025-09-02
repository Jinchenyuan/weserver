package main

import (
	"fmt"
	"net/http"
	"server/core"
	"server/core/config"
	"server/core/transport"
	mesaHttp "server/core/transport/http"
	"time"

	"github.com/gin-gonic/gin"
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
	)

	httpServer := m.GetServerByType(transport.HTTP).(*mesaHttp.Server)
	httpServer.RegisterRoute(http.MethodGet, "/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin hello world"})
	})

	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}
