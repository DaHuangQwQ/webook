package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/ranking/domain"
)

type RedisRankingCache struct {
	redis redis.Cmdable
	key   string
}

func NewRedisRankingCache(redis redis.Cmdable) RankingCache {
	return &RedisRankingCache{
		redis: redis,
		key:   "article:topN",
	}
}

func (cache *RedisRankingCache) Set(ctx context.Context, arts []domain.Article) error {
	for i := 0; i < len(arts); i++ {
		arts[i].Content = ""
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return fmt.Errorf("redis 缓存序列化出错 %v", err)
	}
	// 可以永久缓存
	return cache.redis.Set(ctx, cache.key, val, time.Minute*10).Err()
}

func (cache *RedisRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	bytes, err := cache.redis.Get(ctx, cache.key).Bytes()
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	err = json.Unmarshal(bytes, &arts)
	return arts, err
}
