package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/internal/domain"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserKeyNotExist = redis.Nil
)

type UserCache struct {
	client     redis.Cmdable
	Expiration time.Duration
}

// A 用到了 B, B 一定是 接口
// A 用到了 B, B 一定是 A的字段
// A 用到了 B, A 绝对不初始化 B, 而是外面注入
// 依赖注入

func NewUserCache(redis redis.Cmdable) *UserCache {
	return &UserCache{
		client:     redis,
		Expiration: time.Minute * 15,
	}
}

func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err == nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *UserCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(user.Id)
	return cache.client.Set(ctx, key, val, cache.Expiration).Err()
}

func (cache *UserCache) key(key int64) string {
	return fmt.Sprintf("user:info:%d", key)
}
