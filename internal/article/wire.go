//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/article/events"
	"github.com/DaHuangQwQ/webook/article/grpc"
	"github.com/DaHuangQwQ/webook/article/ioc"
	"github.com/DaHuangQwQ/webook/article/repository"
	"github.com/DaHuangQwQ/webook/article/repository/cache"
	"github.com/DaHuangQwQ/webook/article/repository/dao"
	"github.com/DaHuangQwQ/webook/article/service"
	"github.com/google/wire"
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
