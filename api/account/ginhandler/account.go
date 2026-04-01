package ginhandler

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	pb "server/protobuf/gen"
	"time"

	"github.com/Jinchenyuan/wego"
	"github.com/Jinchenyuan/wego/transport"
	"github.com/Jinchenyuan/wego/transport/micro"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/metadata"
	"go-micro.dev/v5/selector"
)

func getAccountClient() (pb.AccountService, error) {
	m := wego.GetGlobalMesa()
	if m == nil {
		return nil, fmt.Errorf("failed to get global mesa")
	}

	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	clientAny := ms.GetServiceClient("account")
	accountClient, ok := clientAny.(pb.AccountService)
	if !ok {
		return nil, fmt.Errorf("failed to cast to AccountClient")
	}
	return accountClient, nil
}

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
	accountClient, err := getAccountClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	accountClient, err := getAccountClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	rsp, err := accountClient.Login(ctx, &pb.LoginRequest{Account: req.Account, Password: req.Password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(int(rsp.GetCode()), gin.H{
		"code":       rsp.GetCode(),
		"account_id": rsp.GetAccountId(),
		"token":      rsp.GetToken(),
		"message":    rsp.GetMessage(),
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
	accountClient, err := getAccountClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
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
