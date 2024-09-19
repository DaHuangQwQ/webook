package validator

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
	"webook/pkg/logger"
	"webook/pkg/migrator"
	"webook/pkg/migrator/events"
)

type Validator[T migrator.Entity] struct {
	base      *gorm.DB
	target    *gorm.DB
	direction string
	producer  events.Producer
	l         logger.LoggerV1
	batchSize int

	uTime         int64
	sleepInterval time.Duration
	fromBase      func(ctx context.Context, offset int) (T, error)
}

func NewValidator[T migrator.Entity](base *gorm.DB, target *gorm.DB, direction string, l logger.LoggerV1, p events.Producer) *Validator[T] {
	res := &Validator[T]{
		base:          base,
		target:        target,
		direction:     direction,
		producer:      p,
		l:             l,
		batchSize:     100,
		sleepInterval: 0,
	}
	res.fromBase = res.fullFromBase
	return res
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
	offset := 0
	for {
		src, err := v.fromBase(ctx, offset)
		if err == context.DeadlineExceeded || err == context.Canceled {
		}
		if err == gorm.ErrRecordNotFound {
			// 你增量校验，要考虑一直运行的
			// 这个就是咩有数据
			if v.sleepInterval <= 0 {
			}
			time.Sleep(v.sleepInterval)
			continue
		}
		if err != nil {
			// 查询出错了
			v.l.Error("base -> target 查询 base 失败", logger.Error(err))
			// 在这里，
			offset++
			continue
		}

		// 这边就是正常情况
		var dst T
		err = v.target.WithContext(ctx).
			Where("id = ?", src.ID()).
			First(&dst).Error
		switch err {
		case gorm.ErrRecordNotFound:
			// target 没有
			// 丢一条消息到 Kafka 上
			v.notify(src.ID(), events.InconsistentEventTypeTargetMissing)
		case nil:
			equal := src.CompareTo(dst)
			if !equal {
				// 要丢一条消息到 Kafka 上
				v.notify(src.ID(), events.InconsistentEventTypeNotEqual)
			}
		default:
			// 记录日志，然后继续
			// 做好监控
			v.l.Error("base -> target 查询 target 失败",
				logger.Int64("id", src.ID()),
				logger.Error(err))
		}
		offset++
	}
}

func (v *Validator[T]) validateTargetToBase(ctx context.Context) {
	offset := 0
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		var dstTs []T
		err := v.target.WithContext(dbCtx).Select("id").Offset(offset).Limit(v.batchSize).Find(&dstTs).Error
		cancel()
		if len(dstTs) == 0 {
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		}
		switch err {
		case context.Canceled, context.DeadlineExceeded:
			return
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
				v.notifyBaseMissing(diff)
			case gorm.ErrRecordNotFound:
				v.notifyBaseMissing(ids)
			default:
				v.l.Error("查询 base 数据库失败", logger.Error(err))
				continue
			}
		case gorm.ErrRecordNotFound:
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		default:
			v.l.Error("查询 target 数据库失败", logger.Error(err))
			continue
		}
		offset += len(dstTs)
		if len(dstTs) < v.batchSize {
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		}
	}
}

func (v *Validator[T]) Full() *Validator[T] {
	v.fromBase = v.fullFromBase
	return v
}

func (v *Validator[T]) Incr() *Validator[T] {
	v.fromBase = v.incrFromBase
	return v
}

func (v *Validator[T]) UTime(t int64) *Validator[T] {
	v.uTime = t
	return v
}

func (v *Validator[T]) SleepInterval(interval time.Duration) *Validator[T] {
	v.sleepInterval = interval
	return v
}

func (v *Validator[T]) fullFromBase(ctx context.Context, offset int) (T, error) {
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var src T
	err := v.base.WithContext(dbCtx).Order("id").
		Offset(offset).First(&src).Error
	return src, err
}

func (v *Validator[T]) incrFromBase(ctx context.Context, offset int) (T, error) {
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var src T
	err := v.base.WithContext(dbCtx).
		Where("u_time > ?", v.uTime).
		Order("u_time").
		Offset(offset).First(&src).Error
	return src, err
}

func (v *Validator[T]) notifyBaseMissing(ids []int64) {
	for _, id := range ids {
		v.notify(id, events.InconsistentEventTypeBaseMissing)
	}
}

func (v *Validator[T]) notify(id int64, typ string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := v.producer.ProduceInconsistentEvent(ctx, events.InconsistentEvent{
		ID:        id,
		Type:      typ,
		Direction: v.direction,
	})
	if err != nil {
		v.l.Error("发送不一致消息失败",
			logger.Error(err),
			logger.String("type", typ),
			logger.Int64("id", id))
	}
}
