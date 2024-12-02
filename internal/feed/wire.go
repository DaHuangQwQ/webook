//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/feed/grpc"
	"github.com/DaHuangQwQ/webook/feed/ioc"
	"github.com/DaHuangQwQ/webook/feed/repository"
	"github.com/DaHuangQwQ/webook/feed/repository/cache"
	"github.com/DaHuangQwQ/webook/feed/repository/dao"
	"github.com/DaHuangQwQ/webook/feed/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	cache.NewFeedEventCache,
	dao.NewFeedPullEventDAO,
	dao.NewFeedPushEventDAO,
	repository.NewFeedEventRepo,
	service.NewFeedService,
	service.NewFollowEventHandler,
	//service.NewArticleEventHandler,
	//service.NewLikeEventHandler,
	grpc.NewFeedEventGrpcSvc,
)

var thirdSet = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitEtcdClient,
	ioc.NewGrpcxServer,
	ioc.RegisterHandler,
	ioc.InitFollowClient,
)

func initApp() *App {
	wire.Build(
		serverSet,
		thirdSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
