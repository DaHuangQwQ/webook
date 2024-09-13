package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var config Config
	err := viper.UnmarshalKey("redis", &config)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(&redis.Options{
		Addr: config.Addr,
	})
}

func InitRlockClient(client redis.Cmdable) *rlock.Client {
	return rlock.NewClient(client)
}
