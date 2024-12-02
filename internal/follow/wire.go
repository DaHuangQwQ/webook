//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/follow/grpc"
	"github.com/DaHuangQwQ/webook/follow/ioc"
	"github.com/DaHuangQwQ/webook/follow/repository"
	"github.com/DaHuangQwQ/webook/follow/repository/cache"
	"github.com/DaHuangQwQ/webook/follow/repository/dao"
	"github.com/DaHuangQwQ/webook/follow/service"
	"github.com/google/wire"
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
