//go:build wireinject

package ioc

import (
	"github.com/DaHuangQwQ/webook/internal/bff"
	"github.com/google/wire"
)

var thirdSet = wire.NewSet(
	InitRedis,
	//InitDB,
	InitLogger,
)

func InitApp() *App {
	wire.Build(
		thirdSet,

		bff.InitApp,
		wire.FieldsOf(new(*bff.App), "web"),

		//user.InitApp,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
