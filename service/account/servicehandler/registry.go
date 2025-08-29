package servicehandler

import (
	"log"
	pb "server/protobuf/gen"

	"go-micro.dev/v5"
)

func Registry(s micro.Service) error {
	if err := pb.RegisterGreeterHandler(s.Server(), new(Greeter)); err != nil {
		log.Fatalf("register handler: %v", err)
	}
	return nil
}
