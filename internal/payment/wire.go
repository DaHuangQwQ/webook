//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/internal/payment/repository"
	"github.com/DaHuangQwQ/webook/internal/payment/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/payment/web"
	"github.com/google/wire"
)

var ServerSet = wire.NewSet(
	dao.NewPaymentGORMDAO,
	repository.NewPaymentRepository,
	web.NewWechatHandler,
)

func InitApp() *App {
	wire.Build(
		ServerSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
