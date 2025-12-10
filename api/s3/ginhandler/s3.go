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
)

type PutKeyRequest struct {
	Bucket string `json:"bucket" example:"my-bucket"`    // 桶名称
	Key    string `json:"key" example:"images/logo.png"` // 文件路径
	Data   []byte `json:"data"`                          // 文件内容
}

type PutKeyResponse struct {
	Code    int32  `json:"code" example:"200"`                  // 响应码
	Message string `json:"message" example:"PutKey successful"` // 响应消息
}

// PutKey 上传文件
// @Summary 上传文件到 S3
// @Description 上传文件内容
// @Tags S3
// @Accept json
// @Produce json
// @Param request body PutKeyRequest true "上传参数"
// @Success 200 {object} PutKeyResponse
// @Router /s3/PutKey [post]
func PutKey(c *gin.Context) {
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

	rsp, err := s3Client.PutKey(ctx, &pb.PutKeyReq{Key: "test-key", Data: "test-data"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    rsp.GetCode(),
		"message": rsp.GetMessage(),
	})
}

func GetKey(c *gin.Context) {
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

	rsp, err := s3Client.GetKey(ctx, &pb.GetKeyReq{Key: "test-key"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    rsp.GetCode(),
		"message": rsp.GetMessage(),
		"data":    rsp.GetData(),
	})
}

func DeleteKey(c *gin.Context) {
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

	rsp, err := s3Client.DeleteKey(ctx, &pb.DeleteKeyReq{Key: "test-key"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    rsp.GetCode(),
		"message": rsp.GetMessage(),
	})
}
