//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/account/grpc"
	"github.com/DaHuangQwQ/webook/account/ioc"
	"github.com/DaHuangQwQ/webook/account/repository"
	"github.com/DaHuangQwQ/webook/account/repository/dao"
	"github.com/DaHuangQwQ/webook/account/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	dao.NewAccountGORMDAO,
	repository.NewAccountRepository,
	service.NewAccountService,
	grpc.NewAccountServiceServer,
)

var thirdSet = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.NewEtcdClient,
	ioc.NewGrpcServer,
)

func initApp() *App {
	wire.Build(
		serverSet,
		thirdSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
