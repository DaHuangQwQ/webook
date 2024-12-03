//go:build wireinject

package user

import (
	"github.com/DaHuangQwQ/webook/internal/user/repository"
	"github.com/DaHuangQwQ/webook/internal/user/repository/cache"
	"github.com/DaHuangQwQ/webook/internal/user/repository/dao"
	"github.com/DaHuangQwQ/webook/internal/user/service"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var serverSet = wire.NewSet(
	dao.NewUserDao,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService,
)

func InitApp(db *gorm.DB, client redis.Cmdable) *App {
	wire.Build(
		serverSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
