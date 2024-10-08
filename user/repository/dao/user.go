package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"time"
	"webook/bff/api"
)
import "gorm.io/gorm"

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

func (dao *GormUserDao) UpdateStatus(ctx context.Context, user User) error {
	return dao.db.WithContext(ctx).Model(&user).Update("user_status", user.Status).Error
}

func (dao *GormUserDao) DeleteByIds(ctx context.Context, ids []int) error {
	return dao.db.WithContext(ctx).Delete(&User{}, "id in (?)", ids).Error
}

func (dao *GormUserDao) InsertAndGetId(ctx context.Context, user User) (int64, error) {
	err := dao.db.WithContext(ctx).Create(&user).Error
	return user.Id, err
}

func (dao *GormUserDao) FindAll(ctx context.Context, req api.UserSearchReq) (total int, userList []User, err error) {
	var tol int64
	err = dao.db.WithContext(ctx).Model(&User{}).Count(&tol).Error
	err = dao.db.WithContext(ctx).Offset(req.PageNum).Limit(req.PageSize).Order(req.OrderBy).Find(&userList).Error
	return int(tol), userList, err
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
