package fixer

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"webook/pkg/migrator"
	"webook/pkg/migrator/events"
)

type Fixer[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB

	columns []string
}

func NewFixer[T migrator.Entity](base *gorm.DB, target *gorm.DB) (*Fixer[T], error) {
	var t T
	rows, err := target.Model(&t).Limit(1).Rows()
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	return &Fixer[T]{
		base:    base,
		target:  target,
		columns: columns,
	}, nil
}

func (f *Fixer[T]) Fix(ctx context.Context, evt events.InconsistentEvent) error {
	var t T
	err := f.base.WithContext(ctx).Where("id = ?", evt.ID).First(&t).Error
	switch err {
	case nil:
		// base 有数据
		// 可以考虑 WHERE base.u_time >= target.u_time
		return f.target.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(f.columns),
		}).Create(&t).Error
	case gorm.ErrRecordNotFound:
		// base 没数据了
		return f.target.WithContext(ctx).Where("id = ?", evt.ID).Delete(&t).Error
	default:
		return errors.New("未知不一致类型")
	}
}

func (f *Fixer[T]) FixV1(ctx context.Context, evt events.InconsistentEvent) error {
	switch evt.Type {
	case events.InconsistentEventTypeTargetMissing:
		// insert
		var t T
		err := f.base.WithContext(ctx).Where("id = ?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			return nil
		case nil:
			return f.target.WithContext(ctx).Create(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeNotEqual:
		// update
		var t T
		err := f.base.WithContext(ctx).Where("id = ?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			return f.target.WithContext(ctx).Where("id = ?", evt.ID).Delete(&t).Error
		case nil:
			return f.target.WithContext(ctx).Updates(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		// delete
		return f.target.WithContext(ctx).Where("id = ?", evt.ID).Delete(new(T)).Error
	default:
		return errors.New("未知不一致类型")
	}
}
