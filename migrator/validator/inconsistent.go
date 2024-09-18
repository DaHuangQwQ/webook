package validator

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"reflect"
	"time"
	"webook/migrator"
	"webook/migrator/events"
	"webook/pkg/logger"
)

type Validator[T migrator.Migration] struct {
	base      *gorm.DB
	target    *gorm.DB
	p         events.Producer
	l         logger.LoggerV1
	batchSize int
}

func NewValidator[T migrator.Migration](base *gorm.DB, target *gorm.DB, l logger.LoggerV1, p events.Producer) *Validator[T] {
	return &Validator[T]{
		base:      base,
		target:    target,
		p:         p,
		l:         l,
		batchSize: 100,
	}
}

func (v *Validator[T]) Validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		v.validateBaseToTarget(ctx)
		return nil
	})
	eg.Go(func() error {
		v.validateTargetToBase(ctx)
		return nil
	})
	return eg.Wait()
}

func (v *Validator[T]) validateBaseToTarget(ctx context.Context) {
	offset := -1
	for {
		offset++
		var src T
		err := v.base.WithContext(ctx).Offset(offset).First(&src).Error
		switch err {
		case nil:
			var dst T
			err := v.target.WithContext(ctx).Where("id = ?", src.ID()).First(&dst).Error
			switch err {
			case nil:
				if !reflect.DeepEqual(dst, src) {
					v.notify(ctx, src.ID(), events.InconsistentEventTypeNotEqual)
				}
			case gorm.ErrRecordNotFound:
				v.notify(ctx, src.ID(), events.InconsistentEventTypeTargetMissing)
			default:
				v.l.Error("src => dst 查询目标表失败", logger.Error(err))
				continue
			}
		case gorm.ErrRecordNotFound:
			return
		default:
			continue
		}
	}

}

func (v *Validator[T]) validateTargetToBase(ctx context.Context) {
	offset := -v.batchSize
	for {
		offset = offset + v.batchSize
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		var dstTs []T
		err := v.target.WithContext(dbCtx).Select("id").Offset(offset).Limit(v.batchSize).Find(&dstTs).Error
		cancel()
		if len(dstTs) == 0 {
			return
		}
		switch err {
		case nil:
			ids := slice.Map(dstTs, func(idx int, t T) int64 {
				return t.ID()
			})
			var srcTs []T
			err = v.base.WithContext(ctx).Where("id IN ?", ids).Find(&srcTs).Error
			switch err {
			case nil:
				srcIds := slice.Map(srcTs, func(idx int, t T) int64 {
					return t.ID()
				})
				// 计算差集
				diff := slice.DiffSet(ids, srcIds)
				v.notifyBaseMissing(ctx, diff)
			case gorm.ErrRecordNotFound:
				v.notifyBaseMissing(ctx, ids)
			default:
				v.l.Error("查询 base 数据库失败", logger.Error(err))
				continue
			}
		case gorm.ErrRecordNotFound:
			return
		default:
			v.l.Error("查询 target 数据库失败", logger.Error(err))
			continue
		}
		if len(dstTs) < v.batchSize {
			return
		}
	}
}

func (v *Validator[T]) notifyBaseMissing(ctx context.Context, ids []int64) {
	for _, id := range ids {
		v.notify(ctx, id, events.InconsistentEventTypeBaseMissing)
	}
}

func (v *Validator[T]) notify(ctx context.Context, id int64, typ string) {

}
