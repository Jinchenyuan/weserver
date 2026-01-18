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
	"time"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/metadata"
	"go-micro.dev/v5/selector"
)

// Hello 账号服务问候接口
// @Summary 账号服务问候接口
// @Description 向账号服务发送问候请求
// @Tags Account
// @Accept json
// @Produce json
// @Param request body HelloRequest true "问候参数"
// @Success 200 {object} HelloResponse
// @Router /account/hello [get]
func Hello(c *gin.Context) {
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
	rsp, err := accountClient.Hello(ctx, &pb.HelloRequest{Name: "this api account"}, client.WithSelectOption(func(so *selector.SelectOptions) {
		so.Context = ctx
	}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": rsp.GetMessage()})
}

// Login 账号登录接口
// @Summary 账号登录接口
// @Description 用户登录账号
// @Tags Account
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录参数"
// @Success 200 {object} LoginResponse
// @Router /account/login [post]
func Login(c *gin.Context) {
	m := core.GetGlobalMesa()
	if m == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get global mesa"})
		return
	}
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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

	rsp, err := accountClient.Login(ctx, &pb.LoginRequest{Username: req.Username, Password: req.Password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(int(rsp.GetCode()), gin.H{
		"code":    rsp.GetCode(),
		"token":   rsp.GetToken(),
		"message": rsp.GetMessage(),
	})
}

// Register 账号注册接口
// @Summary 账号注册接口
// @Description 用户注册账号
// @Tags Account
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册参数"
// @Success 200 {object} RegisterResponse
// @Router /account/register [post]
func Register(c *gin.Context) {
	m := core.GetGlobalMesa()
	if m == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get global mesa"})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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

	rsp, err := accountClient.Register(ctx, &pb.RegisterRequest{Account: req.Account, Name: req.Name, Password: req.Password, Email: req.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(int(rsp.GetCode()), gin.H{
		"code":    rsp.GetCode(),
		"message": rsp.GetMessage(),
	})
}
