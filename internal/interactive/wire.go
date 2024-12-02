//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/interactive/events"
	"github.com/DaHuangQwQ/webook/interactive/grpc"
	"github.com/DaHuangQwQ/webook/interactive/ioc"
	"github.com/DaHuangQwQ/webook/interactive/repository"
	"github.com/DaHuangQwQ/webook/interactive/repository/cache"
	"github.com/DaHuangQwQ/webook/interactive/repository/dao"
	"github.com/DaHuangQwQ/webook/interactive/service"
	"github.com/google/wire"
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
	ioc.InitEtcdClient,
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
