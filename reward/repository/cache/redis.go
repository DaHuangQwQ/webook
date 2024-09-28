package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/reward/domain"
)

type RewardRedisCache struct {
	client redis.Cmdable
}

func NewRewardRedisCache(client redis.Cmdable) RewardCache {
	return &RewardRedisCache{client: client}
}

func (cache *RewardRedisCache) GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	key := cache.codeURLKey(r)
	data, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.CodeURL{}, err
	}
	var res domain.CodeURL
	err = json.Unmarshal(data, &res)
	return res, err
}

func (cache *RewardRedisCache) CachedCodeURL(ctx context.Context, cu domain.CodeURL, r domain.Reward) error {
	key := cache.codeURLKey(r)
	data, err := json.Marshal(cu)
	if err != nil {
		return err
	}
	// 如果你担心 30 分钟刚好是微信订单过期的问题，那么你可以设置成 29 分钟
	return cache.client.Set(ctx, key, data, time.Minute*30).Err()
}

func (cache *RewardRedisCache) codeURLKey(r domain.Reward) string {
	return fmt.Sprintf("reward:code_url:%s:%d:%d",
		r.Target.Biz, r.Target.BizId, r.Uid)
}
