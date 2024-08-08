package ratelimit

import "context"

type Limiter interface {
	// Limit key 是限流对象， 一般是ip地址
	// Limit 是否触发限流
	// 返回 true，就是触发限流
	Limit(ctx context.Context, key string) (bool, error)
}
