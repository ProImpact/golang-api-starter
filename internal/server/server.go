package server

import (
	"apistarter/internal/config"
	"apistarter/internal/db"
	"apistarter/internal/shutdown"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func NewHttpServer(lc fx.Lifecycle, router *gin.Engine, cfg *config.Configuration, shutDown *shutdown.ShutdownManager) *http.Server {
	s := &http.Server{
		Addr:           cfg.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", s.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", s.Addr)
			go s.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("app shutdown activated")
			var err error
			for _, clean := range shutDown.CleanupFuncs {
				err = clean()
				if err != nil {
					slog.Error(err.Error())
				}
			}
			return s.Shutdown(ctx)
		},
	})
	return s
}

// NewRouter setup the routes here
func NewRouter(q *db.Queries) *gin.Engine {
	e := gin.Default()
	e.GET("/", func(ctx *gin.Context) {
		data, _ := q.GetApplicationName(ctx)
		ctx.JSON(200, data)
	})
	return e
}
