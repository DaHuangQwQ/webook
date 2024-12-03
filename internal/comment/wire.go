//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/internal/comment/repository"
	"github.com/DaHuangQwQ/webook/internal/comment/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/comment/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	dao.NewCommentDAO,
	repository.NewCommentRepo,
	service.NewCommentSvc,
)

func initApp() *App {
	wire.Build(
		serverSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
