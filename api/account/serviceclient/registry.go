package serviceclient

import (
	"fmt"
	"server/core/transport"
	pb "server/protobuf/gen"

	coremicro "server/core/transport/micro"

	"go-micro.dev/v5"
	"go-micro.dev/v5/registry"
)

func Registry(reg registry.Registry) map[string]any {
	ret := make(map[string]any)

	// account
	ret[string(transport.Account)] = make(map[string]any)
	accountService := micro.NewService(
		micro.Name(string(transport.Account)),
		// micro.Registry(reg),
		micro.Selector(coremicro.NewSelectorDependId(reg)),
	)
	accountService.Init()
	ret[fmt.Sprintf("%s-%s", string(transport.Account), "account")] = pb.NewAccountService(string(transport.Account), accountService.Client())

	return ret
}
