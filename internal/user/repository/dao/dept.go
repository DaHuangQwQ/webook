package dao

import (
	"context"
	"gorm.io/gorm"
)

type GormDeptDao struct {
	db *gorm.DB
}

func NewGormDeptDao(db *gorm.DB) DeptDao {
	return &GormDeptDao{
		db: db,
	}
}

func (dao *GormDeptDao) GetList(ctx context.Context, deptName string, status uint, pageNum, pageSize int) ([]SysDept, error) {
	var (
		dataBase = dao.db.WithContext(ctx).Offset(pageNum).Limit(pageSize)
		dept     []SysDept
	)
	if deptName != "" {
		dataBase = dataBase.Where("dept_name like ?", "%"+deptName+"%")
	}
	if status != 0 {
		dataBase = dataBase.Where("status = ?", status)
	}
	err := dataBase.Find(&dept).Error
	return dept, err
}

func (dao *GormDeptDao) Insert(ctx context.Context, dept SysDept) error {
	return dao.db.WithContext(ctx).Create(&dept).Error
}

func (dao *GormDeptDao) Update(ctx context.Context, dept SysDept) error {
	return dao.db.WithContext(ctx).Updates(&dept).Error
}

func (dao *GormDeptDao) Delete(ctx context.Context, dept SysDept) error {
	return dao.db.WithContext(ctx).Delete(&dept).Error
}

func (dao *GormDeptDao) Find(ctx context.Context, dept SysDept) (SysDept, error) {
	var res SysDept
	err := dao.db.WithContext(ctx).Where("id = ?", dept.DeptID).First(&res).Error
	return res, err
}

func (dao *GormDeptDao) FindAll(ctx context.Context) ([]SysDept, error) {
	var res []SysDept
	err := dao.db.WithContext(ctx).Find(&res).Error
	return res, err
}
