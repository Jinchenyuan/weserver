package main

import (
	"fmt"
	"net/http"
	"server/core"
	"server/core/transport"
	mesaHttp "server/core/transport/http"
	"time"

	"github.com/gin-gonic/gin"
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
		core.WithHttpPort(8082),
		core.WithDSN("postgres://user:password@localhost:5432/land_contract?sslmode=disable"),
	)

	httpServer := m.GetServerByType(transport.HTTP).(*mesaHttp.Server)
	httpServer.RegisterRoute(http.MethodGet, "/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin hello world"})
	})

	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}
