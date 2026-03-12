package ginhandler

import (
	"fmt"

	_ "server/api/account/docs"

	"github.com/Jinchenyuan/wego"
	"github.com/Jinchenyuan/wego/transport"
	"github.com/Jinchenyuan/wego/transport/http"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Registry() error {
	m := wego.GetGlobalMesa()
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

func SetAuthMiddleware(authHandler gin.HandlerFunc) error {
	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}

	hs := m.GetServerByType(transport.HTTP).(*http.Server)
	hs.SetAuthMiddleware(authHandler)

	return nil
}
