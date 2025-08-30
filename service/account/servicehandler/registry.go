package servicehandler

import (
	"log"
	pb "server/protobuf/gen"

	"go-micro.dev/v5"
)

func Registry(s micro.Service) error {
	if err := pb.RegisterAccountHandler(s.Server(), new(Account)); err != nil {
		log.Fatalf("register handler: %v", err)
	}
	return nil
}
