//go:build wireinject

package main

import (
	"github.com/DaHuangQwQ/webook/internal/feed/repository"
	"github.com/DaHuangQwQ/webook/internal/feed/repository/cache"
	"github.com/DaHuangQwQ/webook/internal/feed/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/feed/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	cache.NewFeedEventCache,
	dao.NewFeedPullEventDAO,
	dao.NewFeedPushEventDAO,
	repository.NewFeedEventRepo,
	service.NewFeedService,
	service.NewFollowEventHandler,
	//service.NewArticleEventHandler,
	//service.NewLikeEventHandler,
)

func initApp() *App {
	wire.Build(
		serverSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
