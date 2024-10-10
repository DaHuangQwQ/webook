//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/feed/grpc"
	"webook/feed/ioc"
	"webook/feed/repository"
	"webook/feed/repository/cache"
	"webook/feed/repository/dao"
	"webook/feed/service"
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
