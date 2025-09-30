package sql

import "go.uber.org/fx"

var Mod = fx.Options(
	fx.Invoke(MigrateTo),
)
