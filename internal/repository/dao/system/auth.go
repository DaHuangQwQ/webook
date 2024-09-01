package system

import (
	"context"
	"gorm.io/gorm"
)

type AuthDao interface {
	Insert(ctx context.Context, authRule SysAuthRule) error
	Update(ctx context.Context, authRule SysAuthRule) error
	DeleteByIds(ctx context.Context, ids []int64) error
	Find(ctx context.Context, id int64) (SysAuthRule, error)
	FindAll(ctx context.Context) ([]SysAuthRule, error)
}

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

type SysAuthRule struct {
	ID         uint   `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	PID        uint   `gorm:"not null;default:0" json:"pid"`                  // 父ID
	Name       string `gorm:"size:100;not null;default:''" json:"name"`       // 规则名称
	Title      string `gorm:"size:50;not null;default:''" json:"title"`       // 规则名称
	Icon       string `gorm:"size:300;not null;default:''" json:"icon"`       // 图标
	Condition  string `gorm:"size:255;not null;default:''" json:"condition"`  // 条件
	Remark     string `gorm:"size:255;not null;default:''" json:"remark"`     // 备注
	MenuType   uint8  `gorm:"not null;default:0" json:"menu_type"`            // 类型 0目录 1菜单 2按钮
	Weigh      int    `gorm:"not null;default:0" json:"weigh"`                // 权重
	IsHide     uint8  `gorm:"not null;default:0" json:"is_hide"`              // 显示状态
	Path       string `gorm:"size:100;not null;default:''" json:"path"`       // 路由地址
	Component  string `gorm:"size:100;not null;default:''" json:"component"`  // 组件路径
	IsLink     uint8  `gorm:"not null;default:0" json:"is_link"`              // 是否外链 1是 0否
	ModuleType string `gorm:"size:30;not null;default:''" json:"module_type"` // 所属模块
	ModelID    uint   `gorm:"not null;default:0" json:"model_id"`             // 模型ID
	IsIframe   uint8  `gorm:"not null;default:0" json:"is_iframe"`            // 是否内嵌iframe
	IsCached   uint8  `gorm:"not null;default:0" json:"is_cached"`            // 是否缓存
	Redirect   string `gorm:"size:255;not null;default:''" json:"redirect"`   // 路由重定向地址
	IsAffix    uint8  `gorm:"not null;default:0" json:"is_affix"`             // 是否固定
	LinkURL    string `gorm:"size:500;not null;default:''" json:"link_url"`   // 链接地址
	CreatedAt  int64  `gorm:"autoCreateTime" json:"created_at"`               // 创建日期
	UpdatedAt  int64  `gorm:"autoUpdateTime" json:"updated_at"`               // 修改日期
}
