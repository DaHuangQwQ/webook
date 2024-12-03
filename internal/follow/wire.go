//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/internal/follow/repository"
	"github.com/DaHuangQwQ/webook/internal/follow/repository/cache"
	"github.com/DaHuangQwQ/webook/internal/follow/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/follow/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	cache.NewRedisFollowCache,
	dao.NewGORMFollowRelationDAO,
	repository.NewFollowRelationRepository,
	service.NewFollowRelationService,
)

func initApp() *App {
	wire.Build(
		serverSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
