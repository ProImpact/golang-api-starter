package config

import (
	"go.uber.org/zap"
)

func NewLogger(cfg *Configuration) (*zap.Logger, error) {
	if cfg.Mode == "debug" {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
