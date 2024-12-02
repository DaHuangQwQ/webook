package dao

import (
	"context"
	"gorm.io/gorm"
)

type OrderDao interface {
	Create(ctx context.Context, order Order) error
	DeleteByIds(ctx context.Context, ids []int64) error
	FindAll(ctx context.Context) ([]Order, error)
	Update(ctx context.Context, order Order) error
	FindById(ctx context.Context, id uint64) (Order, error)
}

type GormOrderDao struct {
	db *gorm.DB
}

func NewGormOrderDao(db *gorm.DB) OrderDao {
	return &GormOrderDao{
		db: db,
	}
}

func (dao *GormOrderDao) FindById(ctx context.Context, id uint64) (Order, error) {
	var order Order
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	return order, err
}

func (dao *GormOrderDao) Create(ctx context.Context, order Order) error {
	return dao.db.WithContext(ctx).Create(&order).Error
}

func (dao *GormOrderDao) DeleteByIds(ctx context.Context, ids []int64) error {
	return dao.db.WithContext(ctx).Where("id in (?)", ids).Delete(&Order{}).Error
}

func (dao *GormOrderDao) FindAll(ctx context.Context) ([]Order, error) {
	var orders []Order
	err := dao.db.WithContext(ctx).Find(&orders).Error
	return orders, err
}

func (dao *GormOrderDao) Update(ctx context.Context, order Order) error {
	return dao.db.WithContext(ctx).Save(&order).Error
}

type Order struct {
	Id       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `gorm:"type:varchar(255);not null" json:"title"`
	Content  string `gorm:"type:text;not null" json:"content"`
	Author   string `gorm:"type:varchar(255);not null" json:"author"`
	Salary   string `gorm:"type:varchar(255);not null" json:"salary"`
	Duration int64  `gorm:"type:bigint;not null" json:"duration"`
	Address  string `gorm:"type:varchar(255);not null" json:"address"`
	Phone    string `gorm:"type:varchar(255);not null" json:"phone"`
}
