//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/ioc"
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
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,
		ioc.InitOssService,
		ioc.InitCasbinService,
		// DAO 部分
		dao.NewUserDao,
		dao.NewGormArticleDao,
		systemDao.NewGormAuthDao,
		systemDao.NewGormRoleDao,
		systemDao.NewGormDeptDao,

		// cache 部分
		cache.NewCodeCache,
		cache.NewUserCache,
		cache.NewArticleRedisCache,

		// repository 部分
		repository.NewUserRepository,
		repository.NewCodeRepository,
		systemRepository.NewCachedAuthRepository,
		systemRepository.NewCachedRoleRepository,
		systemRepository.NewCachedUserRepository,
		systemRepository.NewCachedDeptRepository,
		repository.NewCachedArticleRepository,

		// Service 部分
		ioc.InitSMSService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		systemService.NewAuthService,
		systemService.NewRoleService,
		systemService.NewSystemService,
		systemService.NewDeptService,

		ioc.InitWechat,

		jwt.NewRedisJWTHandler,
		// handler 部分
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		system.NewAuthHandler,
		system.NewRoleHandler,
		system.NewUserHandler,
		system.NewDeptHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
