package ginhandler

import (
	"fmt"

	"github.com/Jinchenyuan/wego"
	"github.com/Jinchenyuan/wego/transport"
	"github.com/Jinchenyuan/wego/transport/http"
	"github.com/gin-gonic/gin"
)

func Registry() error {
	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	hs := m.GetServerByType(transport.HTTP).(*http.Server)

	hs.RegisterRoute("GET", "/storylines", ListStorylines)
	hs.RegisterRoute("GET", "/storylines/:id", GetStoryline)
	hs.RegisterRoute("POST", "/storylines", CreateStoryline)
	hs.RegisterRoute("PUT", "/storylines/:id", UpdateStoryline)

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
