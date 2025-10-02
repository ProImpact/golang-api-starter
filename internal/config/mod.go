package config

import "go.uber.org/fx"

var Mod = fx.Options(
	fx.Provide(NewConfiguration),
	fx.Provide(NewRedisClient),
	fx.Provide(NewLogger),
)
