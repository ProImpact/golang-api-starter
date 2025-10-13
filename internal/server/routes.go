package server

import (
	"apistarter/internal/config"
	"apistarter/internal/db"
	"apistarter/internal/metrics"
	"apistarter/internal/security"
	"apistarter/internal/server/midleware"
	"apistarter/internal/server/response"
	"apistarter/internal/validation"
	"apistarter/pkg/model"
	"net/http"
	"os"
	"time"

	_ "apistarter/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/otel/trace"
)

var now = time.Now()

func NewRouter(q *db.Queries, cfg *config.Configuration, tracer trace.Tracer) *gin.Engine {
	router := gin.New()
	router.Use(midleware.Recovery())
	router.Use(midleware.RequestID())
	router.Use(metrics.Middleware())
	router.NoRoute(func(ctx *gin.Context) {
		response.Error(ctx, http.StatusNotFound, model.NOT_FOUND, "route not found", map[string]any{
			"method": ctx.Request.Method,
		})
	})
	router.GET("/metrics", metrics.Handler())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(midleware.TracingMiddleware(tracer))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/", func(ctx *gin.Context) {
		data, _ := q.GetApplicationName(ctx)
		response.Success(ctx, data, "user created", nil)
	})
	router.POST("/", func(ctx *gin.Context) {
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
	router.GET("/version", func(ctx *gin.Context) {
		response.Success(ctx, map[string]string{
			"version":    os.Getenv("GIT_TAG"),
			"build_time": os.Getenv("BUILD_TIME"),
			"up_time":    time.Since(now).String(),
		}, "API version", nil)
	})
	router.GET("/token", func(ctx *gin.Context) {
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

	return router
}
