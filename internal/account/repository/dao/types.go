package dao

import (
	"context"
	"gorm.io/gorm"
)

type AccountDAO interface {
	AddActivities(ctx context.Context, activities ...AccountActivity) error
}

// Account 账号本体
// 包含了当前状态
type Account struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	// 我账号是哪个用户的账号
	Uid int64

	// 唯一标识一个账号
	Account int64 `gorm:"uniqueIndex:account_type"`
	Type    uint8 `gorm:"uniqueIndex:account_type"`

	Balance  int64
	Currency string

	UTime int64
	CTime int64
}

// AccountAudit, AccountBank...

type AccountActivity struct {
	Id  int64 `gorm:"primaryKey,autoIncrement"`
	Uid int64

	Biz   string `gorm:"index:biz_type_id"`
	BizId int64  `gorm:"index:biz_type_id"`

	Account     int64 `gorm:"index:account_type"`
	AccountType uint8 `gorm:"index:account_type"`

	// TYPE 入账还是出账
	Amount   int64
	Currency string

	UTime int64
	CTime int64
}

func (AccountActivity) TableName() string {
	return "account_activities"
}

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&Account{}, &AccountActivity{})
}
