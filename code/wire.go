//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/code/grpc"
	"webook/code/ioc"
	"webook/code/repository"
	"webook/code/repository/cache"
	"webook/code/service"
)

var serverSet = wire.NewSet(
	cache.NewCodeCache,
	repository.NewCodeRepository,
	service.NewCodeService,
	grpc.NewCodeServiceServer,
)

var thirdSet = wire.NewSet(
	ioc.InitRedis,
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.NewGrpcxServer,
)

var clientSet = wire.NewSet(
	ioc.InitSmsGRPCClient,
)

func initApp() *App {
	wire.Build(
		serverSet,
		thirdSet,
		clientSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
