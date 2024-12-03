//go:build wireinject

package bff

import (
	"github.com/DaHuangQwQ/gpkg/logger"
	"github.com/DaHuangQwQ/webook/internal/bff/ioc"
	"github.com/DaHuangQwQ/webook/internal/bff/web"
	ijwt "github.com/DaHuangQwQ/webook/internal/bff/web/jwt"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

//go:generate wire
func InitApp(
	l logger.LoggerV1,
	redisClient redis.Cmdable,
) *App {
	wire.Build(
		ioc.InitGinServer,
		ijwt.NewRedisHandler,

		//web.NewArticleHandler,
		web.NewUserHandler,
		//web.NewRewardHandler,
		//jwt.NewRedisHandler,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
