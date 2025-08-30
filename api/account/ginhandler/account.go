package ginhandler

import (
	"context"
	"net/http"
	"server/core"
	"server/core/transport"
	"server/core/transport/micro"
	pb "server/protobuf/gen"
	protocol "server/submodule/protocol/gen/golang"
	"time"

	mgin "server/core/gin"

	"github.com/gin-gonic/gin"
)

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
	clientAny := ms.GetServiceClient(transport.Account, "account")
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

	if err := mgin.WriteResponse(c, rsp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
