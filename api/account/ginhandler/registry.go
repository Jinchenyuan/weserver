package ginhandler

import (
	"fmt"
	"server/core"
	"server/core/transport"
	"server/core/transport/http"

	_ "server/api/account/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Registry() error {
	m := core.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	hs := m.GetServerByType(transport.HTTP).(*http.Server)

	hs.RegisterRoute("POST", "/account/login", Login)
	hs.RegisterRoute("GET", "/account/hello", Hello)
	hs.RegisterRoute("POST", "/account/register", Register)

	hs.GetEngine().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return nil
}
