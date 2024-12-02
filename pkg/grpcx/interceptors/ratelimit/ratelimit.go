package ratelimit

import (
	"context"
	"github.com/DaHuangQwQ/webook/pkg/logger"
	limit "github.com/DaHuangQwQ/webook/pkg/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	limiter limit.Limiter
	key     string
	l       logger.LoggerV1
}

func NewInterceptorBuilder(limiter limit.Limiter, key string, l logger.LoggerV1) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key, l: l}
}

func (i *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		limited, err := i.limiter.Limit(ctx, i.key)
		if err != nil {
			i.l.Error("判断限流出现问题", logger.Error(err))
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		}
		if limited {
			return nil, status.Error(codes.ResourceExhausted, "limit")
		}
		return handler(ctx, req)
	}
}

func (i *InterceptorBuilder) BuildClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		limited, err := i.limiter.Limit(ctx, i.key)
		if err != nil {
			i.l.Error("判断限流出现问题", logger.Error(err))
			return status.Error(codes.ResourceExhausted, err.Error())
		}
		if limited {
			return status.Error(codes.ResourceExhausted, "limit")
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
