package dao

import (
	"context"
	"gorm.io/gorm"
)

type GormAuthDao struct {
	db *gorm.DB
}

func NewGormAuthDao(db *gorm.DB) AuthDao {
	return &GormAuthDao{
		db: db,
	}
}

func (dao *GormAuthDao) Insert(ctx context.Context, authRule SysAuthRule) error {
	return dao.db.WithContext(ctx).Create(&authRule).Error
}

func (dao *GormAuthDao) Update(ctx context.Context, authRule SysAuthRule) error {
	return dao.db.WithContext(ctx).Save(&authRule).Error
}

func (dao *GormAuthDao) DeleteByIds(ctx context.Context, ids []int64) error {
	return dao.db.WithContext(ctx).Where("id in (?)", ids).Delete(&SysAuthRule{}).Error
}

func (dao *GormAuthDao) Find(ctx context.Context, id int64) (SysAuthRule, error) {
	var auth SysAuthRule
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&auth).Error
	return auth, err
}

func (dao GormAuthDao) FindAll(ctx context.Context) ([]SysAuthRule, error) {
	var list []SysAuthRule
	err := dao.db.WithContext(ctx).Find(&list).Error
	return list, err
}
