package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"time"
	"webook/internal/api"
)
import "gorm.io/gorm"

type UserDao interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
	Update(ctx context.Context, user User) error
	FindAll(ctx context.Context, req api.UserSearchReq) (total int, userList []User, err error)
	InsertAndGetId(ctx context.Context, user User) (int64, error)
	DeleteByIds(ctx *gin.Context, ids []int) error
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

func (dao *GormUserDao) DeleteByIds(ctx *gin.Context, ids []int) error {
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

type User struct {
	Id int64 `json:"id" gorm:"primary_key;autoIncrement"`

	WechatOpenId  sql.NullString `json:"wechat_open_id" gorm:"type:varchar(255);unique;"`
	WechatUnionId sql.NullString `json:"wechat_union_id" gorm:"type:varchar(255);unique;"`

	Email      sql.NullString `json:"email" gorm:"type:varchar(100);unique;sql:null"`
	Phone      sql.NullString `json:"phone" gorm:"type:varchar(100);unique;sql:null"`
	Password   string         `json:"password" gorm:"type:varchar(100);not null"`
	Nickname   string         `json:"nickname" gorm:"type:varchar(100)"`
	Birthday   int            `gorm:"type:int(11);not null;default:0;comment:'生日'"`
	UserStatus uint8          `gorm:"type:tinyint unsigned;not null;default:1;comment:'用户状态;0:禁用,1:正常,2:未验证'"`

	Grade  int `json:"grade" gorm:"type:int(11)"`
	Gender int `json:"gender" gorm:"type:int(11)"`

	DeptID      uint64 `gorm:"type:bigint unsigned;not null;default:0;comment:'部门id'"`
	Remark      string `gorm:"type:varchar(255);not null;comment:'备注'"`
	IsAdmin     uint8  `gorm:"type:tinyint;not null;default:1;comment:'是否后台管理员 1 是  0   否'"` // 注意：这里使用了bool类型，根据实际情况可能需要转换为tinyint
	Address     string `gorm:"type:varchar(255);not null;comment:'联系地址'"`
	Describe    string `gorm:"type:varchar(255);not null;comment:'描述信息'"` // 注意：Go中通常使用Description而不是Describe
	LastLoginIP string `gorm:"type:varchar(15);not null;comment:'最后登录ip'"`
	AvatarUrl   string `json:"avatar_url" gorm:"type:varchar(100)"`
	Status      uint   `json:"status"`

	CTime int64 `json:"ctime" gorm:"autoCreateTime:milli"`
	UTime int64 `json:"utime" gorm:"autoUpdateTime:milli"`
}
