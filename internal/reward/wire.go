//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/internal/reward/repository"
	"github.com/DaHuangQwQ/webook/internal/reward/repository/cache"
	"github.com/DaHuangQwQ/webook/internal/reward/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/reward/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	dao.NewRewardGORMDAO,
	cache.NewRewardRedisCache,
	repository.NewRewardRepository,
	service.NewWechatNativeRewardService,
)

func InitApp() *App {
	wire.Build(
		serverSet,
		wire.Struct(new(App), "*"),
	)

	return new(App)
}
