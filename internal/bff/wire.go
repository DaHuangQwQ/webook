//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/bff/ioc"
	"github.com/DaHuangQwQ/webook/bff/web"
	"github.com/DaHuangQwQ/webook/bff/web/jwt"
	"github.com/DaHuangQwQ/webook/pkg/weapp"
	"github.com/google/wire"
)

func InitApp() *weapp.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitRedis,
		ioc.InitEtcdClient,

		web.NewArticleHandler,
		web.NewUserHandler,
		web.NewRewardHandler,
		jwt.NewRedisHandler,

		ioc.InitUserClient,
		ioc.InitIntrClient,
		ioc.InitRewardClient,
		ioc.InitCodeClient,
		ioc.InitArticleClient,
		ioc.InitGinServer,
		wire.Struct(new(weapp.App), "WebServer"),
	)
	return new(weapp.App)
}
