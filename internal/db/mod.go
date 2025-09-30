package db

import "go.uber.org/fx"

var Mod = fx.Options(
	fx.Provide(NewPostgresqlDrive),
	fx.Provide(NewQueries),
)
