package framework

import (
	"crou-api/framework/gormadapter"
	"go.uber.org/fx"
)

var Ctx = fx.Module("framework",
	fx.Provide(
		gormadapter.NewUserGorm,
		gormadapter.NewRoutineGorm,
	),
)
