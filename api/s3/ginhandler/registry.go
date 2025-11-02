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

	hs.RegisterRoute("POST", "/s3/PutKey", PutKey)
	hs.RegisterRoute("POST", "/s3/GetKey", GetKey)

	return nil
}
