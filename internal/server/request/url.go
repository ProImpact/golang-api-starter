package request

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"apistarter/internal/server/response"
	"apistarter/pkg/model"

	"github.com/gin-gonic/gin"
)

func QueryReadString(ctx *gin.Context, key string, defaultValue string) string {
	s := ctx.Query(key)
	if s == "" {
		return defaultValue
	}
	return s
}

// The QueryReadCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found, it returns
// the provided default value.
func QueryReadCSV(ctx *gin.Context, key string, defaultValue []string) []string {
	csv := ctx.Query(key)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

func QueryReadInt(ctx *gin.Context, key string, defaultValue int) int {
	s := ctx.Query(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, model.INVALID_ARGUMENT,
			fmt.Sprintf("invalid query value for %s expected an integer", key),
			nil,
		)
		return -1
	}
	return i
}
