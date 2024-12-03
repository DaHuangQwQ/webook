//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/internal/account/repository"
	"github.com/DaHuangQwQ/webook/internal/account/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/account/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	dao.NewAccountGORMDAO,
	repository.NewAccountRepository,
	service.NewAccountService,
)

func initApp() *App {
	wire.Build(
		serverSet,
		thirdSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
