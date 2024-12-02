//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/tag/events"
	"github.com/DaHuangQwQ/webook/tag/grpc"
	"github.com/DaHuangQwQ/webook/tag/ioc"
	"github.com/DaHuangQwQ/webook/tag/repository"
	"github.com/DaHuangQwQ/webook/tag/repository/cache"
	"github.com/DaHuangQwQ/webook/tag/repository/dao"
	"github.com/DaHuangQwQ/webook/tag/service"
	"github.com/google/wire"
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
