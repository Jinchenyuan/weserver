package ginhandler

import (
	"fmt"
	"server/core"
	"server/core/transport"
	"server/core/transport/http"

	_ "server/api/s3/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Registry() error {
	m := core.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	hs := m.GetServerByType(transport.HTTP).(*http.Server)

	hs.RegisterRoute("POST", "/s3/PutKey", PutKey)
	hs.RegisterRoute("POST", "/s3/GetKey", GetKey)
	hs.RegisterRoute("POST", "/s3/DeleteKey", DeleteKey)

	hs.GetEngine().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return nil
}
