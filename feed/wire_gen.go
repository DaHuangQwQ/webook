// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func initApp() *App {
	loggerV1 := ioc.InitLogger()
	db := ioc.InitDB(loggerV1)
	feedPullEventDAO := dao.NewFeedPullEventDAO(db)
	feedPushEventDAO := dao.NewFeedPushEventDAO(db)
	cmdable := ioc.InitRedis()
	feedEventCache := cache.NewFeedEventCache(cmdable)
	feedEventRepo := repository.NewFeedEventRepo(feedPullEventDAO, feedPushEventDAO, feedEventCache)
	followServiceClient := ioc.InitFollowClient()
	v := ioc.RegisterHandler(feedEventRepo, followServiceClient)
	feedService := service.NewFeedService(feedEventRepo, v)
	feedEventGrpcSvc := grpc.NewFeedEventGrpcSvc(feedService)
	client := ioc.InitEtcdClient()
	server := ioc.NewGrpcxServer(feedEventGrpcSvc, loggerV1, client)
	app := &App{
		server: server,
	}
	return app
}

// wire.go:

var serverSet = wire.NewSet(cache.NewFeedEventCache, dao.NewFeedPullEventDAO, dao.NewFeedPushEventDAO, repository.NewFeedEventRepo, service.NewFeedService, service.NewFollowEventHandler, grpc.NewFeedEventGrpcSvc)

var thirdSet = wire.NewSet(ioc.InitLogger, ioc.InitDB, ioc.InitRedis, ioc.InitEtcdClient, ioc.NewGrpcxServer, ioc.RegisterHandler, ioc.InitFollowClient)
