//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/internal/code/repository"
	"github.com/DaHuangQwQ/webook/internal/code/repository/cache"
	"github.com/DaHuangQwQ/webook/internal/code/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	cache.NewCodeCache,
	repository.NewCodeRepository,
	service.NewCodeService,
)

func initApp() *App {
	wire.Build(
		serverSet,
		thirdSet,
		clientSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
