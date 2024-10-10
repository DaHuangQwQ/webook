//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/tag/events"
	"webook/tag/grpc"
	"webook/tag/ioc"
	"webook/tag/repository"
	"webook/tag/repository/cache"
	"webook/tag/repository/dao"
	"webook/tag/service"
)

var serverSet = wire.NewSet(
	cache.NewRedisTagCache,
	dao.NewGORMTagDAO,
	repository.NewTagRepository,
	service.NewTagService,
	grpc.NewTagServiceServer,
	events.NewSaramaSyncProducer,
)

var thirdSet = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitEtcdClient,
	ioc.NewGrpcxServer,
	ioc.InitKafka,
	ioc.InitSyncProducer,
)

func initApp() *App {
	wire.Build(
		serverSet,
		thirdSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
