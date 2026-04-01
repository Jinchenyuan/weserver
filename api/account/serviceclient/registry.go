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

	// account
	ret["account"] = make(map[string]any)
	accountService := micro.NewService(
		micro.Name("account.client"),
		micro.Selector(coremicro.NewSelectorDependId(reg)),
	)
	accountService.Init()
	ret["account"] = pb.NewAccountService(cfg.Services.Account, accountService.Client())

	return ret
}
