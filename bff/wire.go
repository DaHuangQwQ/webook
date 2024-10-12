//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/bff/ioc"
	"webook/bff/web"
	"webook/bff/web/jwt"
	"webook/pkg/weapp"
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
