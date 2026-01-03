package serviceclient

import (
	"fmt"
	"server/core/config"
	"server/core/transport"
	pb "server/protobuf/gen"

	coremicro "server/core/transport/micro"

	"go-micro.dev/v5"
	"go-micro.dev/v5/registry"
)

func Registry(reg registry.Registry) map[string]any {
	cfg, err := config.Read("config.toml")
	if err != nil {
		panic(fmt.Sprintf("failed to read config: %v", err))
	}
	ret := make(map[string]any)

	// s3
	ret[string(transport.S3)] = make(map[string]any)
	s3Service := micro.NewService(
		micro.Name(string(transport.S3)),
		micro.Selector(coremicro.NewSelectorDependId(reg)),
	)
	s3Service.Init()
	ret[string(transport.S3)] = pb.NewS3Service(cfg.Services.S3, s3Service.Client())

	return ret
}
