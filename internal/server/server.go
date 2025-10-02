package server

import (
	"apistarter/internal/config"
	"apistarter/internal/db"
	"apistarter/internal/server/midleware"
	"apistarter/internal/server/response"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(q *db.Queries, cfg *config.Configuration) *gin.Engine {
	e := gin.Default()
	e.Use(midleware.Recovery())
	e.Use(midleware.RequestID())
	e.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	e.GET("/", func(ctx *gin.Context) {
		data, _ := q.GetApplicationName(ctx)
		response.Success(ctx, data, "user created", nil)
	})
	return e
}
