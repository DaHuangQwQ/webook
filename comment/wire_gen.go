// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"webook/comment/grpc"
	"webook/comment/ioc"
	"webook/comment/repository"
	"webook/comment/repository/dao"
	"webook/comment/service"
)

// Injectors from wire.go:

func initApp() *App {
	loggerV1 := ioc.InitLogger()
	db := ioc.InitDB(loggerV1)
	commentDAO := dao.NewCommentDAO(db)
	commentRepository := repository.NewCommentRepo(commentDAO, loggerV1)
	commentService := service.NewCommentSvc(commentRepository)
	commentServiceServer := grpc.NewCommentServiceServer(commentService)
	client := ioc.InitEtcdClient()
	server := ioc.NewGrpcServer(commentServiceServer, client, loggerV1)
	app := &App{
		server: server,
	}
	return app
}

// wire.go:

var serverSet = wire.NewSet(dao.NewCommentDAO, repository.NewCommentRepo, service.NewCommentSvc, grpc.NewCommentServiceServer)

var thirdSet = wire.NewSet(ioc.InitLogger, ioc.InitDB, ioc.InitEtcdClient, ioc.NewGrpcServer)