package ginhandler

import (
	"fmt"
	"server/core"
	"server/core/transport"
	"server/core/transport/http"
)

func Registry() error {
	m := core.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	hs := m.GetServerByType(transport.HTTP).(*http.Server)

	hs.RegisterRoute("GET", "/account/login", AccountLogin)

	return nil
}
