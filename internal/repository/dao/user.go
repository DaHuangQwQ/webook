package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"time"
)
import "gorm.io/gorm"

type UserDao struct {
	db *gorm.DB
}

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func NewUserDao(db *gorm.DB) *UserDao {
	err := db.AutoMigrate(&User{})
	if err != nil {
		return nil
	}
	return &UserDao{db: db}
}

func (dao *UserDao) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CTime = now
	user.UTime = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error
	return u, err
}

type User struct {
	Id       int64  `json:"id" gorm:"primary_key;autoIncrement"`
	Email    string `json:"email" gorm:"type:varchar(100);not null""'`
	Password string `json:"password" gorm:"type:varchar(100);not null"`

	CTime int64 `json:"ctime" gorm:"autoCreateTime:milli"`
	UTime int64 `json:"utime" gorm:"autoUpdateTime:milli"`
}
