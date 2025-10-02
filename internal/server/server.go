package server

import (
	"apistarter/internal/config"
	"apistarter/internal/db"
	"apistarter/internal/security"
	"apistarter/internal/server/midleware"
	"apistarter/internal/server/response"
	"apistarter/internal/validation"
	"apistarter/pkg/model"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	e.POST("/", func(ctx *gin.Context) {
		var a validation.Address
		err := ctx.ShouldBindBodyWithJSON(&a)
		if err != nil {
			response.Error(
				ctx,
				http.StatusBadRequest,
				model.INVALID_REQUEST,
				"error parsing the request to json",
				map[string]any{
					"error": err.Error(),
				},
			)
			return
		}
		err = a.Validate()
		if err != nil {
			response.Error(
				ctx,
				http.StatusBadRequest,
				model.FIELD_VALIDATION_ERROR,
				"error validating the json payload",
				map[string]any{
					"error": err,
				},
			)
			return
		}
	})
	e.GET("/panic", func(ctx *gin.Context) {
		panic("server recovery")
	})
	e.GET("/token", func(ctx *gin.Context) {
		token, err := security.GenerateToken(uuid.NewString(), time.Minute)
		if err != nil {
			response.Error(
				ctx,
				http.StatusForbidden,
				model.ACCESS_DENIED,
				"error generating the token for the user",
				map[string]any{
					"error": err.Error(),
				},
			)
			return
		}
		response.Success(
			ctx,
			map[string]any{
				"token": token,
			},
			"jwt token generated",
			nil,
		)
	})
	return e
}
