package response

import (
	"apistarter/pkg/model"
	"apistarter/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
			"error":  err.Error(),
			"method": ctx.Request.Method,
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
			"error":  err,
			"method": ctx.Request.Method,
		},
	)
}

func ValidationErrorMsg(ctx *gin.Context, err error) {
	Error(
		ctx,
		http.StatusBadRequest,
		model.FIELD_VALIDATION_ERROR,
		"error validating the json payload",
		map[string]any{
			"error":  err.Error(),
			"method": ctx.Request.Method,
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
			"error":  err,
			"method": ctx.Request.Method,
		},
	)
}

func InternalServerError(ctx *gin.Context, err error) {
	Error(
		ctx,
		http.StatusInternalServerError,
		model.INTERNAL_ERROR,
		"internal server error",
		map[string]any{
			"method": ctx.Request.Method,
		},
	)
}

func InternalServerErrorLog(ctx *gin.Context, l *zap.Logger, err error) {
	InternalServerError(ctx, err)
	l.Error(err.Error())
}
