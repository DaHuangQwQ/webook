package dao

import "context"
import "gorm.io/gorm"

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (u *UserDao) Insert(ctx context.Context, user User) error {
	return u.db.WithContext(ctx).Create(&user).Error
}

type User struct {
	Name     string `gorm:"primaryKey"`
	Age      int
	Password string

	CTime int64
	UTime int64
}
