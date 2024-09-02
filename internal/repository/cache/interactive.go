package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"webook/internal/domain"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const (
	fieldLikeCnt    = "like_cnt"
	fieldReadCnt    = "read_cnt"
	fieldCollectCnt = "collect_cnt"

	ErrKeyNotExist = redis.Nil
)

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, id int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, id int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, id int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, id int64, res domain.Interactive) error
}

type RedisInteractiveCache struct {
	client redis.Cmdable
}

func NewRedisInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{
		client: client,
	}
}

// 缓存方案
// map[string]int

func (cache *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	res, err := cache.client.Eval(ctx, luaIncrCnt,
		[]string{cache.key(biz, bizId)},
		fieldReadCnt, 1).Int()
	if res == 0 {
		// 一般缓存过期
		return errors.New("IncrReadCnt 缓存中 key 不存在")
	}
	return err
}

func (cache *RedisInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, id int64) error {
	res, err := cache.client.Eval(ctx, luaIncrCnt,
		[]string{cache.key(biz, id)},
		fieldLikeCnt, 1).Int()
	if res == 0 {
		// 一般缓存过期
		return errors.New("IncrLikeCnt 缓存中 key 不存在")
	}
	return err
}

func (cache *RedisInteractiveCache) DecrLikeCntIfPresent(ctx context.Context, biz string, id int64) error {
	res, err := cache.client.Eval(ctx, luaIncrCnt,
		[]string{cache.key(biz, id)},
		fieldLikeCnt, -1).Int()
	if res == 0 {
		// 一般缓存过期
		return errors.New("DecrLikeCnt 缓存中 key 不存在")
	}
	return err
}

func (cache *RedisInteractiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, id int64) error {
	res, err := cache.client.Eval(ctx, luaIncrCnt,
		[]string{cache.key(biz, id)},
		fieldCollectCnt, 1).Int()
	if res == 0 {
		// 一般缓存过期
		return errors.New("IncrCollection 缓存中 key 不存在")
	}
	return err
}

func (cache *RedisInteractiveCache) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	key := cache.key(biz, id)
	res, err := cache.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}
	var intr domain.Interactive
	// 这边是可以忽略错误的
	intr.CollectCnt, _ = strconv.ParseInt(res[fieldCollectCnt], 10, 64)
	intr.LikeCnt, _ = strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	intr.ReadCnt, _ = strconv.ParseInt(res[fieldReadCnt], 10, 64)
	return intr, nil
}

func (cache *RedisInteractiveCache) Set(ctx context.Context, biz string, id int64, res domain.Interactive) error {
	key := cache.key(biz, id)
	err := cache.client.HSet(ctx, key, fieldCollectCnt, res.CollectCnt,
		fieldReadCnt, res.ReadCnt,
		fieldLikeCnt, res.LikeCnt,
	).Err()
	if err != nil {
		return err
	}
	return cache.client.Expire(ctx, key, time.Minute*15).Err()
}

func (cache *RedisInteractiveCache) key(biz string, id int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, id)
}
