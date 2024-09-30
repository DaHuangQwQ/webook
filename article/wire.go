//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/article/events"
	"webook/article/grpc"
	"webook/article/ioc"
	"webook/article/repository"
	"webook/article/repository/cache"
	"webook/article/repository/dao"
	"webook/article/service"
)

var serverSet = wire.NewSet(
	dao.NewGormArticleDao,
	cache.NewArticleRedisCache,
	repository.NewCachedArticleRepository,
	repository.NewGrpcAuthorRepository,
	service.NewArticleService,
	grpc.NewArticleServiceServer,
	events.NewKafkaProducer,
)

var thirdSet = wire.NewSet(
	ioc.InitRedis,
	ioc.InitDB,
	ioc.InitEtcdClient,
	ioc.NewGrpcxServer,
	ioc.InitLogger,
	ioc.InitAliyunOss,
	ioc.InitSyncProducer,
	ioc.InitUserRpcClient,
)

func initApp() *App {
	wire.Build(
		serverSet,
		thirdSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
