//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/internal/tag/events"
	"github.com/DaHuangQwQ/webook/internal/tag/repository"
	"github.com/DaHuangQwQ/webook/internal/tag/repository/cache"
	"github.com/DaHuangQwQ/webook/internal/tag/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/tag/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	cache.NewRedisTagCache,
	dao.NewGORMTagDAO,
	repository.NewTagRepository,
	service.NewTagService,
	events.NewSaramaSyncProducer,
)

func initApp() *App {
	wire.Build(
		serverSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
