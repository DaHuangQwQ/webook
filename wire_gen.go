// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"webook/internal/ioc"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/repository/dao/system"
	system2 "webook/internal/repository/system"
	"webook/internal/service"
	system3 "webook/internal/service/system"
	"webook/internal/web"
	"webook/internal/web/jwt"
	system4 "webook/internal/web/system"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	loggerV1 := ioc.InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, loggerV1)
	db := ioc.InitDB(loggerV1)
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmdable)
	client := ioc.InitOssService()
	userRepository := repository.NewUserRepository(userDao, userCache, client)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	codeService := service.NewCodeService(codeRepository, smsService)
	handler := jwt.NewRedisJWTHandler(cmdable)
	userHandler := web.NewUserHandler(userService, codeService, handler, loggerV1)
	wechatService := ioc.InitWechat(loggerV1)
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, handler, userService, loggerV1)
	articleDao := dao.NewGormArticleDao(db)
	articleCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(articleDao, client, userRepository, articleCache)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(articleService, loggerV1)
	authDao := system.NewGormAuthDao(db)
	iEnforcer := ioc.InitCasbinService()
	authRepository := system2.NewCachedAuthRepository(authDao, iEnforcer)
	authService := system3.NewAuthService(authRepository)
	roleDao := system.NewGormRoleDao(db)
	roleRepository := system2.NewCachedRoleRepository(iEnforcer, roleDao)
	roleService := system3.NewRoleService(roleRepository, authRepository)
	authHandler := system4.NewAuthHandler(authService, loggerV1, roleService)
	roleHandler := system4.NewRoleHandler(roleService, loggerV1)
	systemUserRepository := system2.NewCachedUserRepository(iEnforcer)
	systemUserService := system3.NewSystemService(roleService, authService, systemUserRepository)
	systemUserHandler := system4.NewUserHandler(systemUserService, loggerV1)
	deptDao := system.NewGormDeptDao(db)
	deptRepository := system2.NewCachedDeptRepository(deptDao)
	deptService := system3.NewDeptService(deptRepository)
	deptHandler := system4.NewDeptHandler(deptService, loggerV1)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler, articleHandler, authHandler, roleHandler, systemUserHandler, deptHandler)
	return engine
}
