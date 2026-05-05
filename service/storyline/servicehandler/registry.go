package servicehandler

import (
	"log"
	pb "server/protobuf/gen"

	"go-micro.dev/v5"
)

func Registry(s micro.Service) error {
	if err := pb.RegisterStorylineHandler(s.Server(), NewStoryline(nil)); err != nil {
		log.Fatalf("register storyline handler: %v", err)
	}
	return nil
}
