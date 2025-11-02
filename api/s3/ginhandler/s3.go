package ginhandler

import (
	"context"
	"net/http"
	"server/core"
	"server/core/transport"
	"server/core/transport/micro"
	pb "server/protobuf/gen"
	"time"

	"github.com/gin-gonic/gin"

	mgin "server/core/gin"
	protocol "server/submodule/protocol/gen/golang"
)

func PutKey(c *gin.Context) {
	putKeyReq := &protocol.S3PutKeyReq{}
	if err := mgin.ReadRequest(c, putKeyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := core.GetGlobalMesa()
	if m == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get global mesa"})
		return
	}

	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	clientAny := ms.GetServiceClient(transport.S3)
	s3Client, ok := clientAny.(pb.S3Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get s3 client"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rsp, err := s3Client.PutKey(ctx, &pb.PutKeyReq{Key: putKeyReq.Key, Data: putKeyReq.Data})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpRsp := &protocol.S3PutKeyResp{
		Code:    protocol.ErrorCode(rsp.GetCode()),
		Message: rsp.GetMessage(),
	}
	if err := mgin.WriteResponse(c, httpRsp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func GetKey(c *gin.Context) {
	getKeyReq := &protocol.S3GetKeyReq{}
	if err := mgin.ReadRequest(c, getKeyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := core.GetGlobalMesa()
	if m == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get global mesa"})
		return
	}

	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	clientAny := ms.GetServiceClient(transport.S3)
	s3Client, ok := clientAny.(pb.S3Service)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get s3 client"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rsp, err := s3Client.GetKey(ctx, &pb.GetKeyReq{Key: getKeyReq.Key})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpRsp := &protocol.S3GetKeyResp{
		Code:    protocol.ErrorCode(rsp.GetCode()),
		Data:    rsp.GetData(),
		Message: rsp.GetMessage(),
	}
	if err := mgin.WriteResponse(c, httpRsp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
