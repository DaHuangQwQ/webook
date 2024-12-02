package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type GormRoleDao struct {
	db *gorm.DB
}

func NewGormRoleDao(db *gorm.DB) RoleDao {
	return &GormRoleDao{
		db: db,
	}
}

func (g *GormRoleDao) GetRoleListSearch(ctx context.Context, role SysRole, pageNum int, pageSize int) (int64, []SysRole, error) {
	var (
		total    int64
		dateBase = g.db.WithContext(ctx).Offset(pageNum - 1).Limit(pageSize)
		roles    []SysRole
	)
	err := g.db.WithContext(ctx).Model(SysRole{}).Count(&total).Error
	if err != nil {
		return 0, nil, fmt.Errorf("count role list error: %w", err)
	}
	if role.Name != "" {
		dateBase = dateBase.Where("name like ?", "%"+role.Name+"%")
	}
	err = dateBase.Find(&roles).Error
	return total, roles, err
}

func (g *GormRoleDao) FindById(ctx context.Context, roleId int64) (SysRole, error) {
	var role SysRole
	err := g.db.WithContext(ctx).Where("id = ?", roleId).First(&role).Error
	return role, err
}

func (g *GormRoleDao) Create(ctx context.Context, role SysRole) error {
	return g.db.WithContext(ctx).Create(&role).Error
}

func (g *GormRoleDao) CreateAndGetId(ctx context.Context, role SysRole) (id int64, err error) {
	err = g.db.WithContext(ctx).Create(&role).Error
	return role.ID, err
}

func (g *GormRoleDao) UpdateById(ctx context.Context, role SysRole) error {
	return g.db.WithContext(ctx).Save(role).Error
}

func (g *GormRoleDao) FindList(ctx context.Context) ([]SysRole, error) {
	var list []SysRole
	err := g.db.WithContext(ctx).Find(&list).Error
	return list, err
}

func (g *GormRoleDao) DeleteByIds(ctx context.Context, roleIds []int64) error {
	return g.db.WithContext(ctx).Where("id in (?)", roleIds).Delete(&SysRole{}).Error
}
