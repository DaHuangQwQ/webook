//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/search/grpc"
	"webook/search/ioc"
	"webook/search/repository"
	"webook/search/repository/dao"
	"webook/search/service"
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
