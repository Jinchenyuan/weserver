package serviceclient

import (
	"fmt"
	pb "server/protobuf/gen"

	coremicro "github.com/Jinchenyuan/wego/transport/micro"

	"github.com/Jinchenyuan/wego/config"
	"github.com/Jinchenyuan/wego/transport"
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
	ret[string(transport.Account)] = make(map[string]any)
	accountService := micro.NewService(
		micro.Name(string(transport.Account)),
		micro.Selector(coremicro.NewSelectorDependId(reg)),
	)
	accountService.Init()
	ret[string(transport.Account)] = pb.NewAccountService(cfg.Services.Account, accountService.Client())

	return ret
}
