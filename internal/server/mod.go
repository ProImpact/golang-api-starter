package server

import "go.uber.org/fx"

var Mod = fx.Options(
	fx.Provide(NewRouter),
	fx.Provide(NewHttpServer),
)
