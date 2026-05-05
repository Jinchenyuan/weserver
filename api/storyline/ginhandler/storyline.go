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

func getStorylineClient() (pb.StorylineService, error) {
	m := wego.GetGlobalMesa()
	if m == nil {
		return nil, errors.New("failed to get global mesa")
	}
	ms := m.GetServerByType(transport.MICRO_SERVER).(*micro.Service)
	clientAny := ms.GetServiceClient("storyline")
	storylineClient, ok := clientAny.(pb.StorylineService)
	if !ok {
		return nil, errors.New("failed to cast to StorylineService")
	}
	return storylineClient, nil
}

func ListStorylines(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getStorylineClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rsp, err := client.ListStorylines(ctx, &pb.ListStorylinesRequest{AccountId: accountID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toListResponse(rsp))
}

func GetStoryline(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getStorylineClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rsp, err := client.GetStoryline(ctx, &pb.GetStorylineRequest{
		AccountId: accountID,
		Id:        c.Param("id"),
	})
	if err != nil {
		if errors.Is(err, model.ErrStorylineNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "storyline not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toDetail(rsp))
}

func CreateStoryline(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getStorylineClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req CreateStorylineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validateNodes(req.Nodes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rsp, err := client.CreateStoryline(ctx, toCreateRequest(accountID, req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	statusCode := http.StatusCreated
	if rsp.GetCode() != 0 {
		statusCode = int(rsp.GetCode())
	}
	c.JSON(statusCode, toMutationResponse(rsp))
}

func UpdateStoryline(c *gin.Context) {
	accountID, ok := commonmiddleware.GetAccountID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing account context"})
		return
	}
	client, err := getStorylineClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req UpdateStorylineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.ID != "" && req.ID != c.Param("id") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path id and body id do not match"})
		return
	}
	req.ID = c.Param("id")
	if err := validateNodes(req.Nodes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rsp, err := client.UpdateStoryline(ctx, toUpdateRequest(accountID, req))
	if err != nil {
		if errors.Is(err, model.ErrStorylineNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "storyline not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	statusCode := http.StatusOK
	if rsp.GetCode() != 0 {
		statusCode = int(rsp.GetCode())
	}
	c.JSON(statusCode, toMutationResponse(rsp))
}
