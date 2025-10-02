package app

import (
	"apistarter/internal/server/midleware"
	"apistarter/pkg/model"
	"time"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, code model.ErrorCode, message string, details map[string]any) {
	fault := "client"
	if status >= 500 {
		fault = "server"
	}
	err := model.RequestErr{
		Code:      code,
		Message:   message,
		Details:   details,
		TimeStamp: time.Now().UTC(),
		Path:      c.Request.URL.Path,
		RequestId: midleware.GetRequestID(c),
		Status:    status,
		Fault:     fault,
	}
	c.JSON(status, err)
}
