//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/user/grpc"
	"webook/user/ioc"
	"webook/user/repository"
	"webook/user/repository/cache"
	"webook/user/repository/dao"
	system2 "webook/user/repository/system"
	"webook/user/service"
	"webook/user/service/system"
	"webook/user/web"
)

var serverSet = wire.NewSet(
	dao.NewUserDao,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService,
	grpc.NewUserServiceServer,
)

var thirdSet = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitOssService,
	ioc.InitCasbinService,
	ioc.InitEtcdClient,
	ioc.NewGrpcxServer,
)

func initApp() *App {
	wire.Build(
		serverSet,
		thirdSet,

		dao.NewGormDeptDao,
		dao.NewGormAuthDao,
		dao.NewGormRoleDao,

		system2.NewCachedAuthRepository,
		system2.NewCachedDeptRepository,
		system2.NewCachedUserRepository,
		system2.NewCachedRoleRepository,

		system.NewAuthService,
		system.NewSystemService,
		system.NewDeptService,
		system.NewRoleService,
		system.NewSysMonitorService,

		web.NewRoleHandler,
		web.NewUserHandler,
		web.NewMonitorHandler,
		web.NewDeptHandler,
		web.NewAuthHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
