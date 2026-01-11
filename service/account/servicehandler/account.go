package servicehandler

import (
	"context"
	"fmt"
	"server/core"
	"server/model"
	pb "server/protobuf/gen"
	"server/utils"

	"github.com/google/uuid"
)

type Account struct{}

func (a *Account) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	// Implement your login logic here
	fmt.Printf("Login request received: username=%s, password=%s\n", req.GetUsername(), req.GetPassword())

	rsp.Code = 200
	rsp.Token = "some-token"
	rsp.Message = "Login successful"
	return nil
}

func (a *Account) Hello(ctx context.Context, req *pb.HelloRequest, rsp *pb.HelloResponse) error {
	fmt.Printf("Hello request received: name=%s\n", req.GetName())
	rsp.Message = "Hello, " + req.GetName()
	return nil
}

func (a *Account) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.RegisterResponse) error {
	// Implement your registration logic here
	fmt.Printf("Register request received: account=%s, password=%s, email=%s\n", req.GetAccount(), req.GetPassword(), req.GetEmail())

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
