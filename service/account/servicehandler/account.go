package servicehandler

import (
	"context"
	"fmt"
	pb "server/protobuf/gen"
)

type Account struct{}

func (a *Account) Login(ctx context.Context, req *pb.AccountLoginReq, rsp *pb.AccountLoginResp) error {
	// Implement your login logic here
	fmt.Printf("Login request received: username=%s, password=%s\n", req.GetUsername(), req.GetPassword())

	rsp.Code = 200
	rsp.Token = "some-token"
	rsp.Message = "Login successful"
	return nil
}

func (a *Account) Hello(ctx context.Context, req *pb.AccountHelloReq, rsp *pb.AccountHelloResp) error {
	fmt.Printf("Hello request received: name=%s\n", req.GetName())
	rsp.Message = "Hello, " + req.GetName()
	return nil
}

func (a *Account) Register(ctx context.Context, req *pb.AccountRegisterReq, rsp *pb.AccountRegisterResp) error {
	// Implement your registration logic here
	fmt.Printf("Register request received: username=%s, password=%s, email=%s\n", req.GetUsername(), req.GetPassword(), req.GetEmail())

	rsp.Code = 201
	rsp.Message = "Registration successful"
	return nil
}
