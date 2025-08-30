package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type request struct {
	Message []byte `json:"message" binding:"required"`
}

func ReadRequest(c *gin.Context, r protoreflect.ProtoMessage) error {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		return err
	}
	if err := proto.Unmarshal(req.Message, r); err != nil {
		return err
	}

	return nil
}

func WriteResponse(c *gin.Context, r protoreflect.ProtoMessage) error {
	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, gin.H{"message": data})
	return nil
}
