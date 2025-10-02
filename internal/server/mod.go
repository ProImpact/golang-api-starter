package server

import (
	"apistarter/internal/config"
	"apistarter/internal/shutdown"
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Mod = fx.Options(
	fx.Provide(NewRouter),
	fx.Provide(NewHttpServer),
)

func NewHttpServer(
	lc fx.Lifecycle,
	router *gin.Engine,
	cfg *config.Configuration, shutDown *shutdown.ShutdownManager,
	logger *zap.Logger,
) *http.Server {
	// add instrumentation to the hole app
	handler := otelhttp.NewHandler(router, "/")
	s := &http.Server{
		Addr:           cfg.Port,
		Handler:        handler,
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
			logger.Info("Starting HTTP server at", zap.String("port", s.Addr))
			go s.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("app shutdown activated")
			var err error
			for _, clean := range shutDown.CleanupFuncs {
				err = clean()
				if err != nil {
					slog.Error(err.Error())
				}
			}
			for _, clean := range shutDown.CleanupFuncsWithContext {
				err = clean(ctx)
				if err != nil {
					slog.Error(err.Error())
				}
			}
			return s.Shutdown(ctx)
		},
	})
	return s
}
