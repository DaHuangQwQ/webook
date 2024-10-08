//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/follow/grpc"
	"webook/follow/ioc"
	"webook/follow/repository"
	"webook/follow/repository/cache"
	"webook/follow/repository/dao"
	"webook/follow/service"
)

var serverSet = wire.NewSet(
	cache.NewRedisFollowCache,
	dao.NewGORMFollowRelationDAO,
	repository.NewFollowRelationRepository,
	service.NewFollowRelationService,
	grpc.NewFollowServiceServer,
)

var thirdSet = wire.NewSet(
	ioc.InitRedis,
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitEtcdClient,
	ioc.NewGrpcServer,
)

func initApp() *App {
	wire.Build(
		serverSet,
		thirdSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
