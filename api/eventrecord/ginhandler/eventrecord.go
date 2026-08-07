package ginhandler

import (
	"context"
	"errors"
	"net/http"
	commonmiddleware "server/api/middleware"
	"server/model"
	pb "server/protobuf/gen"
	"time"

	"github.com/Jinchenyuan/wego"
	"github.com/Jinchenyuan/wego/transport"
	"github.com/Jinchenyuan/wego/transport/micro"
	"github.com/gin-gonic/gin"
)

func getEventRecordClient() (pb.EventRecordService, error) {
	m := wego.GetGlobalMesa()
	if m == nil {
		return nil, errors.New("failed to get global mesa")
	}
	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	clientAny := ms.GetServiceClient("eventrecord")
	eventClient, ok := clientAny.(pb.EventRecordService)
	if !ok {
		return nil, errors.New("failed to cast to EventRecordService")
	}
	return eventClient, nil
}

func ListEventRecords(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getEventRecordClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rsp, err := client.ListEventRecords(ctx, &pb.ListEventRecordsRequest{AccountId: accountID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toEventListResponse(rsp))
}

func GetEventRecord(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getEventRecordClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rsp, err := client.GetEventRecord(ctx, &pb.GetEventRecordRequest{
		AccountId: accountID,
		Id:        c.Param("id"),
	})
	if err != nil {
		if errors.Is(err, model.ErrEventRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toEventDetail(rsp))
}

func CreateEventRecord(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getEventRecordClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req CreateEventRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validateEventRecordRequest(req.Title, req.OccurredAt, req.CoverPhotoURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rsp, err := client.CreateEventRecord(ctx, toCreateEventRecordRequest(accountID, req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	statusCode := http.StatusCreated
	if rsp.GetCode() != 0 {
		statusCode = int(rsp.GetCode())
	}
	c.JSON(statusCode, toEventMutationResponse(rsp))
}

func UpdateEventRecord(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getEventRecordClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req UpdateEventRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.ID != "" && req.ID != c.Param("id") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path id and body id do not match"})
		return
	}
	req.ID = c.Param("id")
	if err := validateEventRecordRequest(req.Title, req.OccurredAt, req.CoverPhotoURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rsp, err := client.UpdateEventRecord(ctx, toUpdateEventRecordRequest(accountID, req))
	if err != nil {
		if errors.Is(err, model.ErrEventRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	statusCode := http.StatusOK
	if rsp.GetCode() != 0 {
		statusCode = int(rsp.GetCode())
	}
	c.JSON(statusCode, toEventMutationResponse(rsp))
}

func CreateEventNote(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getEventRecordClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req CreateEventNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validateEventNoteRequest(req.Title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rsp, err := client.CreateEventNote(ctx, toCreateEventNoteRpcRequest(accountID, c.Param("id"), req))
	if err != nil {
		if errors.Is(err, model.ErrEventRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	statusCode := http.StatusCreated
	if rsp.GetCode() != 0 {
		statusCode = int(rsp.GetCode())
	}
	c.JSON(statusCode, toEventMutationResponse(rsp))
}

func UpdateEventNote(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getEventRecordClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req UpdateEventNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.EventID != "" && req.EventID != c.Param("eventId") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path eventId and body eventId do not match"})
		return
	}
	if req.NoteID != "" && req.NoteID != c.Param("noteId") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path noteId and body noteId do not match"})
		return
	}
	req.EventID = c.Param("eventId")
	req.NoteID = c.Param("noteId")
	if err := validateEventNoteRequest(req.Title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rsp, err := client.UpdateEventNote(ctx, toUpdateEventNoteRpcRequest(accountID, req))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEventRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "event record not found"})
		case errors.Is(err, model.ErrEventNoteNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "event note not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	statusCode := http.StatusOK
	if rsp.GetCode() != 0 {
		statusCode = int(rsp.GetCode())
	}
	c.JSON(statusCode, toEventMutationResponse(rsp))
}

func DeleteEventNote(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getEventRecordClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rsp, err := client.DeleteEventNote(ctx, &pb.DeleteEventNoteRequest{
		AccountId: accountID,
		EventId:   c.Param("eventId"),
		NoteId:    c.Param("noteId"),
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEventRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "event record not found"})
		case errors.Is(err, model.ErrEventNoteNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "event note not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	statusCode := http.StatusOK
	if rsp.GetCode() != 0 {
		statusCode = int(rsp.GetCode())
	}
	c.JSON(statusCode, toEventMutationResponse(rsp))
}
