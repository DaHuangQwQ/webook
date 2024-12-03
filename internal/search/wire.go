//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/internal/search/repository"
	"github.com/DaHuangQwQ/webook/internal/search/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/search/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	dao.NewArticleElasticDAO,
	dao.NewAnyESDAO,
	dao.NewTagESDAO,
	dao.NewUserElasticDAO,
	repository.NewAnyRepository,
	repository.NewArticleRepository,
	repository.NewUserRepository,
	service.NewSearchService,
	service.NewSyncService,
)

func initApp() *App {
	wire.Build(
		serverSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
