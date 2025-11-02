package ginhandler

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"server/core"
	"server/core/transport"
	"server/core/transport/micro"
	pb "server/protobuf/gen"
	protocol "server/submodule/protocol/gen/golang"
	"time"

	mgin "server/core/gin"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/metadata"
	"go-micro.dev/v5/selector"
)

func AccountHello(c *gin.Context) {
	m := core.GetGlobalMesa()
	if m == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get global mesa"})
		return
	}

	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	clientAny := ms.GetServiceClient(transport.Account)
	accountClient, ok := clientAny.(pb.AccountService)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cast to AccountClient"})
		return
	}
	// 随机生成一个6位数的uid
	uid := rand.Intn(900000) + 100000
	ctx := metadata.NewContext(context.Background(), map[string]string{"uid": fmt.Sprintf("%d", uid)})
	// For testing, use a fixed ui
	// ctx := metadata.NewContext(context.Background(), map[string]string{"uid": "123456"})
	rsp, err := accountClient.Hello(ctx, &pb.AccountHelloReq{Name: "this api account"}, client.WithSelectOption(func(so *selector.SelectOptions) {
		so.Context = ctx
	}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": rsp.GetMessage()})
}

func AccountLogin(c *gin.Context) {
	loginReq := &protocol.AccountLoginReq{}
	if err := mgin.ReadRequest(c, loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := core.GetGlobalMesa()
	if m == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get global mesa"})
		return
	}

	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	clientAny := ms.GetServiceClient(transport.Account)
	accountClient, ok := clientAny.(pb.AccountService)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cast to AccountClient"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rsp, err := accountClient.Login(ctx, &pb.AccountLoginReq{Username: loginReq.GetUsername(), Password: loginReq.GetPassword()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpRsp := &protocol.AccountLoginResp{
		Code:    protocol.ErrorCode(rsp.GetCode()),
		Token:   rsp.GetToken(),
		Message: rsp.GetMessage(),
	}
	if err := mgin.WriteResponse(c, httpRsp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
