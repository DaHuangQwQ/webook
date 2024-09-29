//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/article/repository"
	"webook/article/repository/cache"
	"webook/article/repository/dao"
	"webook/article/service"
)

var serverSet = wire.NewSet(
	dao.NewGormArticleDao,
	cache.NewArticleRedisCache,
	repository.NewCachedArticleRepository,
	service.NewArticleService,
)

var thirdSet = wire.NewSet()

func initApp() *App {
	wire.Build(
		serverSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
