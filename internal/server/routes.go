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
	"time"

	_ "apistarter/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewRouter(q *db.Queries, cfg *config.Configuration, tracer trace.Tracer) *gin.Engine {
	router := gin.New()
	router.Use(midleware.Recovery())
	router.Use(midleware.RequestID())
	router.GET("/metrics", metrics.Handler())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
	router.GET("/panic", func(ctx *gin.Context) {
		_, span := tracer.Start(ctx.Request.Context(), "panic")
		defer span.End()

		userID := "123"

		span.SetAttributes(
			attribute.String("http.method", ctx.Request.Method),
			attribute.String("http.route", "/panic"),
			attribute.String("user.id", userID),
		)
		panic("random error")
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
