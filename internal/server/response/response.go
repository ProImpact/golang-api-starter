package response

import (
	"apistarter/pkg/model"
	"apistarter/pkg/utils"
	"net/http"
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
		RequestId: utils.GetRequestID(c),
		Status:    status,
		Fault:     fault,
	}
	c.JSON(status, err)
}

func Success(c *gin.Context, data any, message string, meta map[string]any) {
	resp := model.Success{
		Data:      data,
		Message:   message,
		Meta:      meta,
		RequestId: utils.GetRequestID(c),
		TimeStamp: time.Now().UTC(),
	}
	c.JSON(http.StatusOK, resp)
}
