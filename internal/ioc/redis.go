package ioc

import (
	"github.com/redis/go-redis/v9"
	"webook/config"
)

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Host,
	})
}
