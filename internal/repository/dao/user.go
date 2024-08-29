package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"time"
)
import "gorm.io/gorm"

type UserDao interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
	Update(ctx context.Context, user User) error
}

type GormUserDao struct {
	db *gorm.DB
}

var (
	ErrDuplicate      = errors.New("邮箱或手机冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func NewUserDao(db *gorm.DB) UserDao {
	return &GormUserDao{db: db}
}

func (dao *GormUserDao) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CTime = now
	user.UTime = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicate
		}
	}
	return err
}

func (dao *GormUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *GormUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&u).Error
	return u, err
}

func (dao *GormUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error
	return u, err
}

func (dao *GormUserDao) FindByWechat(ctx context.Context, openId string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id=?", openId).First(&u).Error
	return u, err
}

func (dao *GormUserDao) Update(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.UTime = now
	err := dao.db.WithContext(ctx).Updates(user).Error
	return err
}

type User struct {
	Id int64 `json:"id" gorm:"primary_key;autoIncrement"`

	WechatOpenId  sql.NullString `json:"wechat_open_id" gorm:"type:varchar(255);unique;"`
	WechatUnionId sql.NullString `json:"wechat_union_id" gorm:"type:varchar(255);unique;"`

	Email    sql.NullString `json:"email" gorm:"type:varchar(100);unique;sql:null"`
	Phone    sql.NullString `json:"phone" gorm:"type:varchar(100);unique;sql:null"`
	Password string         `json:"password" gorm:"type:varchar(100);not null"`
	Nickname string         `json:"nickname" gorm:"type:varchar(100)"`
	Grade    int            `json:"grade" gorm:"type:int(11)"`
	Gender   int            `json:"gender" gorm:"type:int(11)"`

	AvatarUrl string `json:"avatar_url" gorm:"type:varchar(100)"`

	CTime int64 `json:"ctime" gorm:"autoCreateTime:milli"`
	UTime int64 `json:"utime" gorm:"autoUpdateTime:milli"`
}
