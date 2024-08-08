package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"webook/internal/service/sms"
)

type FailOverSMSService struct {
	svcs []sms.Service

	// v1 的字段
	// 当前服务商下标
	idx uint64
}

func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

func (f *FailOverSMSService) Send(ctx context.Context, tplId string, args []sms.NamedArg, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("轮询了所有的服务商，但是发送都失败了")
}

// 起始下标轮询
// 并且出错也轮询
func (f *FailOverSMSService) SendV1(ctx context.Context, tplId string, args []sms.NamedArg, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	// 我要迭代 length
	for i := idx; i < idx+length; i++ {
		// 取余数来计算下标
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			// 前者是被取消，后者是超时
			return err
		}
		log.Println(err)
	}
	return errors.New("轮询了所有的服务商，但是发送都失败了")
}
