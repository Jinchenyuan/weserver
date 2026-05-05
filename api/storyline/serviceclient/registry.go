package serviceclient

import (
	"fmt"
	"server/config"
	pb "server/protobuf/gen"

	coremicro "github.com/Jinchenyuan/wego/transport/micro"
	"go-micro.dev/v5"
	"go-micro.dev/v5/registry"
)

func Registry(reg registry.Registry) map[string]any {
	cfg, err := config.Read("config.toml")
	if err != nil {
		panic(fmt.Sprintf("failed to read config: %v", err))
	}
	ret := make(map[string]any)

	storylineService := micro.NewService(
		micro.Name("storyline.client"),
		micro.Selector(coremicro.NewSelectorDependId(reg)),
	)
	storylineService.Init()
	ret["storyline"] = pb.NewStorylineService(cfg.Services.Storyline, storylineService.Client())

	return ret
}
