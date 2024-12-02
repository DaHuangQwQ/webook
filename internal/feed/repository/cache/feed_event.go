package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var FolloweesNotFound = redis.Nil

type FeedEventCache interface {
	SetFollowees(ctx context.Context, follower int64, followees []int64) error
	GetFollowees(ctx context.Context, follower int64) ([]int64, error)
}

type feedEventCache struct {
	client redis.Cmdable
}

func NewFeedEventCache(client redis.Cmdable) FeedEventCache {
	return &feedEventCache{
		client: client,
	}
}

const FolloweeKeyExpiration = 10 * time.Minute

func (f *feedEventCache) SetFollowees(ctx context.Context, follower int64, followees []int64) error {
	key := f.getFolloweeKey(follower)
	followeesStr, err := json.Marshal(followees)
	if err != nil {
		return err
	}
	return f.client.Set(ctx, key, followeesStr, FolloweeKeyExpiration).Err()
}

func (f *feedEventCache) GetFollowees(ctx context.Context, follower int64) ([]int64, error) {
	key := f.getFolloweeKey(follower)
	res, err := f.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, FolloweesNotFound
	}
	var followees []int64
	err = json.Unmarshal([]byte(res), &followees)
	if err != nil {
		return nil, err
	}
	return followees, nil
}

func (f *feedEventCache) getFolloweeKey(follower int64) string {
	return fmt.Sprintf("feed_event:%d", follower)
}
