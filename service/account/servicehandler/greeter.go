package servicehandler

import (
	"context"
	pb "server/protobuf/gen"
)

type Greeter struct{}

func (g *Greeter) SayHello(ctx context.Context, req *pb.HelloRequest, rsp *pb.HelloReply) error {
	rsp.Message = "hello, " + req.GetName() + " this is service."
	return nil
}

func (g *Greeter) SayHelloAgain(ctx context.Context, req *pb.HelloRequest, rsp *pb.HelloReply) error {
	rsp.Message = "hello again, " + req.GetName() + " this is service."
	return nil
}
