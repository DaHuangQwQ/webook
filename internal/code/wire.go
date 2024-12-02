//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/code/grpc"
	"github.com/DaHuangQwQ/webook/code/ioc"
	"github.com/DaHuangQwQ/webook/code/repository"
	"github.com/DaHuangQwQ/webook/code/repository/cache"
	"github.com/DaHuangQwQ/webook/code/service"
	"github.com/google/wire"
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
