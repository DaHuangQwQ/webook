package validator

import (
	"context"
	"gorm.io/gorm"
	"webook/pkg/logger"
	"webook/pkg/migrator"
	events2 "webook/pkg/migrator/events"
)

type CanalIncrValidator[T migrator.Entity] struct {
	baseValidator
}

func NewCanalIncrValidator[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
	direction string,
	l logger.LoggerV1,
	producer events2.Producer,
) *CanalIncrValidator[T] {
	return &CanalIncrValidator[T]{
		baseValidator: baseValidator{
			base:      base,
			target:    target,
			direction: direction,
			l:         l,
			producer:  producer,
		},
	}
}

// Validate 一次校验一条
func (v *CanalIncrValidator[T]) Validate(ctx context.Context, id int64) error {
	var base T

	err := v.base.WithContext(ctx).Where("id = ?").First(&base).Error
	switch err {
	case gorm.ErrRecordNotFound:
		var target T
		err1 := v.target.WithContext(ctx).Where("id = ?").First(&target).Error
		switch err1 {
		case gorm.ErrRecordNotFound:
			// 数据一致
		case nil:
			v.notify(id, events2.InconsistentEventTypeBaseMissing)
		default:
			return err
		}
	case nil:
		var target T
		err1 := v.target.WithContext(ctx).Where("id = ?").First(&target).Error
		switch err1 {
		case gorm.ErrRecordNotFound:
			v.notify(id, events2.InconsistentEventTypeTargetMissing)
		case nil:
			if !base.CompareTo(target) {
				v.notify(id, events2.InconsistentEventTypeNotEqual)
			}
		default:
			return err
		}
	default:
		return err
	}
	return nil
}
