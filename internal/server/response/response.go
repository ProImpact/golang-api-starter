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

func InvalidJsonPayload(ctx *gin.Context, err error) {
	Error(
		ctx,
		http.StatusBadRequest,
		model.INVALID_REQUEST,
		"error parsing the request to json",
		map[string]any{
			"error": err.Error(),
		},
	)
}

func ValidationError(ctx *gin.Context, err error) {
	Error(
		ctx,
		http.StatusBadRequest,
		model.FIELD_VALIDATION_ERROR,
		"error validating the json payload",
		map[string]any{
			"error": err,
		},
	)
}

func QueryValidationError(ctx *gin.Context, err error) {
	Error(
		ctx,
		http.StatusBadRequest,
		model.FIELD_VALIDATION_ERROR,
		"error validating the query params",
		map[string]any{
			"error": err,
		},
	)
}
