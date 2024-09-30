//go:build wireinject

package main

import (
	"github.com/google/wire"
	events2 "webook/article/events"
	repository3 "webook/code/repository"
	cache3 "webook/code/repository/cache"
	service2 "webook/code/service"
	"webook/interactive/events"
	repository2 "webook/interactive/repository"
	cache2 "webook/interactive/repository/cache"
	dao2 "webook/interactive/repository/dao"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/jwt"
	"webook/internal/web/system"
	"webook/ioc"
	repository4 "webook/user/repository"
	"webook/user/repository/cache"
	dao3 "webook/user/repository/dao"
	system4 "webook/user/repository/dao/system"
	system3 "webook/user/repository/system"
	service3 "webook/user/service"
	system2 "webook/user/service/system"
)

func InitWebServer() *App {
	wire.Build(
		// 第三方依赖
		ioc.InitEtcd,
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,
		ioc.InitOssService,
		ioc.InitCasbinService,
		ioc.InitConsumers,
		ioc.InitSyncProducer,
		ioc.InitSaramaClient,
		ioc.InitRlockClient,
		ioc.InitIntrGRPCClient,

		ioc.InitJobs,
		ioc.InitRankingJob,

		// events
		events2.NewKafkaProducer,
		events.NewInteractiveReadEventConsumer,
		// DAO 部分
		dao3.NewUserDao,
		system4.NewGormAuthDao,
		system4.NewGormRoleDao,
		system4.NewGormDeptDao,
		dao2.NewGormInteractiveDao,
		dao.NewGormOrderDao,
		dao.NewGormRecruitmentDao,

		// cache 部分
		cache3.NewCodeCache,
		cache.NewUserCache,
		cache2.NewRedisInteractiveCache,

		// repository 部分
		repository4.NewUserRepository,
		repository3.NewCodeRepository,
		system3.NewCachedAuthRepository,
		system3.NewCachedRoleRepository,
		system3.NewCachedUserRepository,
		system3.NewCachedDeptRepository,
		repository2.NewCachedInteractiveRepository,
		repository.NewCachedOrderRepository,
		repository.NewCachedRecruitmentRepository,

		// Service 部分
		ioc.InitSMSService,
		service3.NewUserService,
		service2.NewCodeService,
		service.NewOrderService,
		service.NewRecruitmentService,
		system2.NewAuthService,
		system2.NewRoleService,
		system2.NewSystemService,
		system2.NewDeptService,
		system2.NewSysMonitorService,

		ioc.InitWechat,

		jwt.NewRedisJWTHandler,
		// handler 部分
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		web.NewOrderHandler,
		web.NewRecruitmentHandler,
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
