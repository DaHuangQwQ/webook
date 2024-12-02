package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type AccountGORMDAO struct {
	db *gorm.DB
}

func NewAccountGORMDAO(db *gorm.DB) AccountDAO {
	return &AccountGORMDAO{db: db}
}

func (dao *AccountGORMDAO) AddActivities(ctx context.Context, activities ...AccountActivity) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		for _, activity := range activities {
			err := dao.db.Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]any{
					"balance": gorm.Expr("`balance`+?", activity.Amount),
					"u_time":  now,
				}),
			}).Create(&Account{
				Uid:      activity.Uid,
				Account:  activity.Account,
				Type:     activity.AccountType,
				Balance:  activity.Amount,
				Currency: activity.Currency,
				UTime:    now,
				CTime:    now,
			}).Error
			if err != nil {
				return err
			}
		}
		return tx.Create(&activities).Error
	})
}
