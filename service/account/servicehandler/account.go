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
