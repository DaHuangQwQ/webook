package circuitbreaker

import (
	"context"
	"github.com/go-kratos/aegis/circuitbreaker"
	"google.golang.org/grpc"
)

type InterceptorBuilder struct {
	breaker circuitbreaker.CircuitBreaker
}

func (b *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		err := b.breaker.Allow()
		if err == nil {
			res, er := handler(ctx, req)
			// 借助这个判断业务错误
			//s, ok := status.FromError(er)
			if er != nil {
				// 没有区别业务错误和系统错误
				b.breaker.MarkFailed()
			} else {
				b.breaker.MarkSuccess()
			}
			return res, err
		}
		b.breaker.MarkFailed()
		return nil, err
	}
}
