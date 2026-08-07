package servicehandler

import (
	"log"
	pb "server/protobuf/gen"

	"go-micro.dev/v5"
)

func Registry(s micro.Service) error {
	if err := pb.RegisterEventRecordHandler(s.Server(), NewEventRecord(nil)); err != nil {
		log.Fatalf("register eventrecord handler: %v", err)
	}
	return nil
}
