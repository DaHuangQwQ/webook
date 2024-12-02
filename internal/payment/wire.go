//go:build wireinject

package main

import (
	grpc2 "github.com/DaHuangQwQ/webook/payment/grpc"
	"github.com/DaHuangQwQ/webook/payment/ioc"
	"github.com/DaHuangQwQ/webook/payment/repository"
	"github.com/DaHuangQwQ/webook/payment/repository/dao"
	"github.com/DaHuangQwQ/webook/payment/web"
	"github.com/google/wire"
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
