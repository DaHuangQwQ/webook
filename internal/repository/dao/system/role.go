package system

import (
	"context"
	"gorm.io/gorm"
)

type RoleDao interface {
	FindList(ctx context.Context) ([]SysRole, error)
	UpdateById(ctx context.Context, role SysRole) error
	Create(ctx context.Context, role SysRole) error
	FindById(ctx context.Context, roleId int64) (SysRole, error)
	DeleteByIds(ctx context.Context, roleIds []int64) error
	CreateAndGetId(ctx context.Context, role SysRole) (id int64, err error)
	GetRoleListSearch(ctx context.Context, role SysRole, pageNum int, pageSize int) ([]SysRole, error)
}

type GormRoleDao struct {
	db *gorm.DB
}

func NewGormRoleDao(db *gorm.DB) RoleDao {
	return &GormRoleDao{
		db: db,
	}
}

func (g *GormRoleDao) GetRoleListSearch(ctx context.Context, role SysRole, pageNum int, pageSize int) ([]SysRole, error) {
	var (
		dateBase = g.db.WithContext(ctx).Offset(pageNum).Limit(pageSize)
		roles    []SysRole
	)
	if role.Name != "" {
		dateBase = dateBase.Where("name like ?", "%"+role.Name+"%")
	}
	err := dateBase.Find(&roles).Error
	return roles, err
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

type SysRole struct {
	ID        int64  `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Status    uint8  `gorm:"type:tinyint(3) unsigned;not null;default:0;comment:状态;0:禁用;1:正常" json:"status"`
	ListOrder uint   `gorm:"type:int(10) unsigned;not null;default:0;comment:排序" json:"list_order"`
	Name      string `gorm:"type:varchar(20);not null;default:'';comment:角色名称" json:"name"`
	Remark    string `gorm:"type:varchar(255);not null;default:'';comment:备注" json:"remark"`
	DataScope uint8  `gorm:"type:tinyint(3) unsigned;not null;default:3;comment:数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）" json:"data_scope"`
	CTime     int64  `gorm:"index;comment:创建时间" json:"c_time"`
	UTime     int64  `gorm:"index;comment:更新时间" json:"u_time"`
}
