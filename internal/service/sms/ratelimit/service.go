package ratelimit

import (
	"context"
	"errors"
	"webook/internal/service/sms"
	limiter "webook/pkg/ratelimit"
)

// 装饰器模式
var errLimited = errors.New("触发限流")

var _ sms.Service = &RateLimitSMSService{}

type RateLimitSMSService struct {
	// 被装饰的
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

// 通过组合实现 装饰器模式
type RateLimitSMSServiceV1 struct {
	sms.Service
	limiter limiter.Limiter
	key     string
}

func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []sms.NamedArg, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errLimited
	}
	return r.svc.Send(ctx, tplId, args, numbers...)
}

func NewRateLimitSMSService(svc sms.Service,
	l limiter.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: l,
		key:     "sms-limiter",
	}
}
