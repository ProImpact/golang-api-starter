package shutdown

import "go.uber.org/fx"

var Mod = fx.Options(
	fx.Provide(NewShutDownManager),
)
