//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"webook/interactive/grpc"
	"webook/interactive/repository"
	"webook/interactive/repository/cache"
	"webook/interactive/repository/dao"
	"webook/interactive/service"
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
