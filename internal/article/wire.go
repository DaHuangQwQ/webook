//go:build wireinject

package article

import (
	"github.com/DaHuangQwQ/webook/internal/article/events"
	"github.com/DaHuangQwQ/webook/internal/article/repository"
	"github.com/DaHuangQwQ/webook/internal/article/repository/cache"
	"github.com/DaHuangQwQ/webook/internal/article/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/article/service"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	dao.NewGormArticleDao,
	cache.NewArticleRedisCache,
	repository.NewCachedArticleRepository,
	repository.NewGrpcAuthorRepository,
	service.NewArticleService,
	events.NewKafkaProducer,
)

func initApp() *App {
	wire.Build(
		serverSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
