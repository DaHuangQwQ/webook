//go:build wireinject

package main

import (
	grpc2 "github.com/DaHuangQwQ/webook/reward/grpc"
	"github.com/DaHuangQwQ/webook/reward/ioc"
	"github.com/DaHuangQwQ/webook/reward/repository"
	"github.com/DaHuangQwQ/webook/reward/repository/cache"
	"github.com/DaHuangQwQ/webook/reward/repository/dao"
	"github.com/DaHuangQwQ/webook/reward/service"
	"github.com/google/wire"
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
