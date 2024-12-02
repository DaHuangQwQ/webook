//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/comment/grpc"
	"github.com/DaHuangQwQ/webook/comment/ioc"
	"github.com/DaHuangQwQ/webook/comment/repository"
	"github.com/DaHuangQwQ/webook/comment/repository/dao"
	"github.com/DaHuangQwQ/webook/comment/service"
	"github.com/google/wire"
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
