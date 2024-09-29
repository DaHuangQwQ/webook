//go:build wireinject

package main

import (
	"github.com/google/wire"
	grpc2 "webook/reward/grpc"
	"webook/reward/ioc"
	"webook/reward/repository"
	"webook/reward/repository/cache"
	"webook/reward/repository/dao"
	"webook/reward/service"
)

var serverSet = wire.NewSet(
	dao.NewRewardGORMDAO,
	cache.NewRewardRedisCache,
	repository.NewRewardRepository,
	service.NewWechatNativeRewardService,
	grpc2.NewRewardServiceServer,
)

var thirdSet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.NewEtcdClient,
	ioc.NewGrpcServer,
	ioc.InitPaymentGrpcClient,
	ioc.InitAccountGrpcClient,
)

func InitApp() *App {
	wire.Build(
		serverSet,
		thirdSet,
		wire.Struct(new(App), "*"),
	)

	return new(App)
}
