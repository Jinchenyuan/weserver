package servicehandler

import (
	"context"
	"fmt"
	"server/model"
	pb "server/protobuf/gen"
	"server/utils"
	"time"

	"github.com/Jinchenyuan/wego/core"
	"github.com/Jinchenyuan/wego/core/logger"
	"github.com/google/uuid"
)

type Account struct {
	log *logger.Logger
}

func NewAccount(log *logger.Logger) *Account {
	return &Account{log: resolveLogger(log)}
}

func resolveLogger(log *logger.Logger) *logger.Logger {
	if log != nil {
		return log
	}

	if globalLog := core.GetGlobalLogger(); globalLog != nil {
		return globalLog
	}

	return logger.GetLogger("account.service")
}

func (a *Account) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	a.log.Info("Login request received: account=", req.GetAccount())
	m := core.GetGlobalMesa()
	if m == nil {
		rsp.Code = 500
		rsp.Message = "failed to get global mesa"
		return nil
	}

	account, err := model.FindAccountByAccount(ctx, m.DB, req.GetAccount())
	if err != nil || account == nil {
		rsp.Code = 401
		rsp.Message = "invalid username or password"
		return nil
	}

	if !utils.CheckPassword(account.Password, req.GetPassword()) {
		rsp.Code = 401
		rsp.Message = "invalid username or password"
		return nil
	}

	token, err := utils.GenerateToken(account.ID)
	if err != nil {
		rsp.Code = 500
		rsp.Message = fmt.Sprintf("failed to generate token: %v", err)
		return nil
	}

	m.Redis.Set(ctx, fmt.Sprintf("token:%d", account.ID), token, time.Hour*24*7)

	rsp.Code = 200
	rsp.Token = token
	rsp.AccountId = account.ID
	rsp.Message = "Login successful"
	return nil
}

func (a *Account) Hello(ctx context.Context, req *pb.HelloRequest, rsp *pb.HelloResponse) error {
	a.log.Info("Hello request received: name=", req.GetName())
	rsp.Message = "Hello, " + req.GetName()
	return nil
}

func (a *Account) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.RegisterResponse) error {
	a.log.Info("Register request received: account=", req.GetAccount(), "email=", req.GetEmail())

	m := core.GetGlobalMesa()
	if m == nil {
		rsp.Code = 500
		rsp.Message = "failed to get global mesa"
		return nil
	}

	existAccount, _ := model.FindAccountByAccount(ctx, m.DB, req.Account)
	if existAccount != nil {
		rsp.Code = 400
		rsp.Message = "account already exists"
		return nil
	}

	hashPasswd, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		rsp.Code = 500
		rsp.Message = fmt.Sprintf("failed to hash password: %v", err)
		return nil
	}

	account := &model.Account{
		ID:       uuid.New().ID(),
		Account:  req.GetAccount(),
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: hashPasswd,
	}
	account.SetDB(m.DB)

	err = account.Create(ctx)
	if err != nil {
		rsp.Code = 500
		rsp.Message = fmt.Sprintf("failed to put key: %v", err)
		return nil
	}

	rsp.Code = 201
	rsp.Message = "Registration successful"
	return nil
}
