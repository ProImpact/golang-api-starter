package midleware

import (
	"apistarter/pkg/model"
	"apistarter/pkg/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		err := model.RequestErr{
			Code:      model.INTERNAL_ERROR,
			Message:   "An unexpected error occurred",
			Details:   map[string]any{"error": fmt.Sprintf("%v", recovered)},
			TimeStamp: time.Now(),
			Path:      c.Request.URL.Path,
			RequestId: utils.GetRequestID(c),
			Status:    http.StatusInternalServerError,
			Fault:     "server",
		}
		c.JSON(http.StatusInternalServerError, err)
	})
}
