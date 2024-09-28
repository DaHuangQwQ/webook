//go:build wireinject

package main

import (
	"github.com/google/wire"
	grpc2 "webook/payment/grpc"
	"webook/payment/ioc"
	"webook/payment/repository"
	"webook/payment/repository/dao"
	"webook/payment/web"
)

var ServerSet = wire.NewSet(
	dao.NewPaymentGORMDAO,
	repository.NewPaymentRepository,
	web.NewWechatHandler,
	grpc2.NewPaymentServiceServer,
)

var thirdSet = wire.NewSet(
	ioc.InitLogger,
	ioc.NewEtcdClient,
	ioc.InitDB,
	ioc.NewGrpcServer,
)

var wechatSet = wire.NewSet(
	ioc.InitWechatConfig,
	ioc.InitWechatNotifyHandler,
	ioc.InitWechatClient,
	ioc.InitWechatNativeService,
)

func InitApp() *App {
	wire.Build(
		ServerSet,
		thirdSet,
		wechatSet,
		ioc.NewWebServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
