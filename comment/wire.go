//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/comment/grpc"
	"webook/comment/ioc"
	"webook/comment/repository"
	"webook/comment/repository/dao"
	"webook/comment/service"
)

var serverSet = wire.NewSet(
	dao.NewCommentDAO,
	repository.NewCommentRepo,
	service.NewCommentSvc,
	grpc.NewCommentServiceServer,
)

var thirdSet = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
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
