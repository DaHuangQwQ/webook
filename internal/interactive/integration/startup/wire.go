//go:build wireinject

package startup

import (
	"github.com/DaHuangQwQ/webook/interactive/grpc"
	"github.com/DaHuangQwQ/webook/interactive/repository"
	"github.com/DaHuangQwQ/webook/interactive/repository/cache"
	"github.com/DaHuangQwQ/webook/interactive/repository/dao"
	"github.com/DaHuangQwQ/webook/interactive/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	InitRedis, InitTestDB,
	InitLog,
	InitKafka,
)

func InitGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(
		grpc.NewInteractiveServiceServer,
		thirdProvider,
		dao.NewGormInteractiveDao,
		cache.NewRedisInteractiveCache,
		repository.NewCachedInteractiveRepository,
		service.NewInteractiveService,
	)
	return new(grpc.InteractiveServiceServer)
}
