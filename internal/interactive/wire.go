//go:build wireinject

package interactive

import (
	"github.com/DaHuangQwQ/webook/internal/interactive/events"
	"github.com/DaHuangQwQ/webook/internal/interactive/repository"
	"github.com/DaHuangQwQ/webook/internal/interactive/repository/cache"
	"github.com/DaHuangQwQ/webook/internal/interactive/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/interactive/service"
	"github.com/google/wire"
)

var interactiveServerProviderSet = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	cache.NewRedisInteractiveCache,
	dao.NewGormInteractiveDao,
)

func InitAPP() *App {
	wire.Build(
		interactiveServerProviderSet,
		events.NewInteractiveReadEventConsumer,
		//ioc.InitConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
