//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/user/grpc"
	"github.com/DaHuangQwQ/webook/user/ioc"
	"github.com/DaHuangQwQ/webook/user/repository"
	"github.com/DaHuangQwQ/webook/user/repository/cache"
	"github.com/DaHuangQwQ/webook/user/repository/dao"
	system2 "github.com/DaHuangQwQ/webook/user/repository/system"
	"github.com/DaHuangQwQ/webook/user/service"
	"github.com/DaHuangQwQ/webook/user/service/system"
	"github.com/DaHuangQwQ/webook/user/web"
	"github.com/google/wire"
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
