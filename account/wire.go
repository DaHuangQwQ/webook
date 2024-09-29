//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/account/grpc"
	"webook/account/ioc"
	"webook/account/repository"
	"webook/account/repository/dao"
	"webook/account/service"
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
