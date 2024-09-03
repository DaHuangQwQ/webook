//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/internal/events/article"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	systemDao "webook/internal/repository/dao/system"
	systemRepository "webook/internal/repository/system"
	"webook/internal/service"
	systemService "webook/internal/service/system"
	"webook/internal/web"
	"webook/internal/web/jwt"
	"webook/internal/web/system"
	"webook/ioc"
)

func InitWebServer() *App {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,
		ioc.InitOssService,
		ioc.InitCasbinService,
		ioc.InitConsumers,
		ioc.InitSyncProducer,
		ioc.InitSaramaClient,
		// events
		article.NewKafkaProducer,
		article.NewInteractiveReadEventConsumer,
		// DAO 部分
		dao.NewUserDao,
		dao.NewGormArticleDao,
		systemDao.NewGormAuthDao,
		systemDao.NewGormRoleDao,
		systemDao.NewGormDeptDao,
		dao.NewGormInteractiveDao,
		dao.NewGormOrderDao,

		// cache 部分
		cache.NewCodeCache,
		cache.NewUserCache,
		cache.NewArticleRedisCache,
		cache.NewRedisInteractiveCache,

		// repository 部分
		repository.NewUserRepository,
		repository.NewCodeRepository,
		systemRepository.NewCachedAuthRepository,
		systemRepository.NewCachedRoleRepository,
		systemRepository.NewCachedUserRepository,
		systemRepository.NewCachedDeptRepository,
		repository.NewCachedArticleRepository,
		repository.NewCachedInteractiveRepository,
		repository.NewCachedOrderRepository,

		// Service 部分
		ioc.InitSMSService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		service.NewInteractiveService,
		service.NewOrderService,
		systemService.NewAuthService,
		systemService.NewRoleService,
		systemService.NewSystemService,
		systemService.NewDeptService,
		systemService.NewSysMonitorService,

		ioc.InitWechat,

		jwt.NewRedisJWTHandler,
		// handler 部分
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		web.NewOrderHandler,
		system.NewAuthHandler,
		system.NewRoleHandler,
		system.NewUserHandler,
		system.NewDeptHandler,
		system.NewMonitorHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
