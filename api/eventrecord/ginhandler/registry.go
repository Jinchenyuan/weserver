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

	hs.RegisterRoute("GET", "/event-records", ListEventRecords)
	hs.RegisterRoute("GET", "/event-records/:id", GetEventRecord)
	hs.RegisterRoute("POST", "/event-records", CreateEventRecord)
	hs.RegisterRoute("PUT", "/event-records/:id", UpdateEventRecord)
	hs.RegisterRoute("POST", "/event-records/:id/notes", CreateEventNote)
	hs.RegisterRoute("PUT", "/event-records/:eventId/notes/:noteId", UpdateEventNote)
	hs.RegisterRoute("DELETE", "/event-records/:eventId/notes/:noteId", DeleteEventNote)

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
