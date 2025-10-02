package server

import (
	"apistarter/internal/db"
	"apistarter/internal/server/midleware"
	"apistarter/internal/server/response"

	"github.com/gin-gonic/gin"
)

func NewRouter(q *db.Queries) *gin.Engine {
	e := gin.Default()
	e.Use(midleware.Recovery())
	e.Use(midleware.RequestID())
	e.GET("/", func(ctx *gin.Context) {
		data, _ := q.GetApplicationName(ctx)
		response.Success(ctx, data, "User created", nil)
	})
	return e
}
