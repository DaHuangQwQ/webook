//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/search/grpc"
	"github.com/DaHuangQwQ/webook/search/ioc"
	"github.com/DaHuangQwQ/webook/search/repository"
	"github.com/DaHuangQwQ/webook/search/repository/dao"
	"github.com/DaHuangQwQ/webook/search/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	dao.NewArticleElasticDAO,
	dao.NewAnyESDAO,
	dao.NewTagESDAO,
	dao.NewUserElasticDAO,
	repository.NewAnyRepository,
	repository.NewArticleRepository,
	repository.NewUserRepository,
	service.NewSearchService,
	service.NewSyncService,
	grpc.NewSearchSyncServiceServer,
	grpc.NewSearchServiceServer,
)

var thirdSet = wire.NewSet(
	ioc.InitLogger,
	ioc.InitESClient,
	ioc.InitGRPCxServer,
	ioc.InitEtcdClient,
)

func initApp() *App {
	wire.Build(
		serverSet,
		thirdSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
