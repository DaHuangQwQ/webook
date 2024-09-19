//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/interactive/events"
	"webook/interactive/grpc"
	"webook/interactive/ioc"
	"webook/interactive/repository"
	"webook/interactive/repository/cache"
	"webook/interactive/repository/dao"
	"webook/interactive/service"
)

var interactiveServerProviderSet = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	cache.NewRedisInteractiveCache,
	dao.NewGormInteractiveDao,
)

var thirdPartySet = wire.NewSet(
	// ioc.InitDB,
	ioc.InitDST,
	ioc.InitSRC,
	ioc.InitBizDB,
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitKafka,
	//ioc.InitSaramaClient,
	ioc.InitDoubleWritePool,
	ioc.InitSyncProducer,
)

var migratorProviderSet = wire.NewSet(
	ioc.InitMigratorWeb,
	ioc.InitMigradatorProducer,
	ioc.InitFixDataConsumer,
	ioc.NewConsumers,
)

func InitAPP() *App {
	wire.Build(
		thirdPartySet,
		interactiveServerProviderSet,
		migratorProviderSet,
		events.NewInteractiveReadEventConsumer,
		grpc.NewInteractiveServiceServer,
		ioc.NewGrpcxServer,
		//ioc.InitConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
