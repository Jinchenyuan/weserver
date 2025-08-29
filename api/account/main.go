package main

import (
	"context"
	"fmt"
	"net/http"
	"server/api/account/serviceclient"
	"server/core"
	"server/core/transport"
	mesaHttp "server/core/transport/http"
	"server/core/transport/micro"
	pb "server/protobuf/gen"
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
		core.WithHttpPort(8083),
		core.WithDSN("postgres://user:password@localhost:5432/land_contract?sslmode=disable"),
	)

	hs := m.GetServerByType(transport.HTTP).(*mesaHttp.Server)
	hs.RegisterRoute(http.MethodGet, "/account", func(c *gin.Context) {
		ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
		clientAny := ms.GetServiceClient(transport.Account, "greeter")
		greeterClient, ok := clientAny.(pb.GreeterService)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cast to GreeterClient"})
			return
		}
		rsp, err := greeterClient.SayHello(context.Background(), &pb.HelloRequest{Name: "api account"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": rsp.Message})
	})

	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	ms.NewServiceClients(serviceclient.Registry)

	if err := m.Run(); err != nil {
		fmt.Printf("failed to run mesa: %v\n", err)
	}
}
