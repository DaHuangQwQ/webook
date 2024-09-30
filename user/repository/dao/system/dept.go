package system

import (
	"context"
	"gorm.io/gorm"
)

type DeptDao interface {
	Insert(ctx context.Context, dept SysDept) error
	Update(ctx context.Context, dept SysDept) error
	Delete(ctx context.Context, dept SysDept) error
	Find(ctx context.Context, dept SysDept) (SysDept, error)
	FindAll(ctx context.Context) ([]SysDept, error)
	GetList(ctx context.Context, deptName string, status uint, pageNum, pageSize int) ([]SysDept, error)
}

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

type SysDept struct {
	DeptID    int64  `gorm:"primaryKey;autoIncrement;column:dept_id" json:"dept_id"`        // 部门id
	ParentID  int64  `gorm:"column:parent_id;default:0" json:"parent_id"`                   // 父部门id
	Ancestors string `gorm:"type:varchar(50);column:ancestors;default:''" json:"ancestors"` // 祖级列表
	DeptName  string `gorm:"type:varchar(30);column:dept_name;default:''" json:"dept_name"` // 部门名称
	OrderNum  int    `gorm:"column:order_num;default:0" json:"order_num"`                   // 显示顺序
	Leader    string `gorm:"type:varchar(20);column:leader" json:"leader"`                  // 负责人
	Phone     string `gorm:"type:varchar(11);column:phone" json:"phone"`                    // 联系电话
	Email     string `gorm:"type:varchar(50);column:email" json:"email"`                    // 邮箱
	Status    uint8  `gorm:"column:status;default:0" json:"status"`                         // 部门状态（0正常 1停用）
	CreatedBy uint64 `gorm:"column:created_by;default:0" json:"created_by"`                 // 创建人
	UpdatedBy int64  `gorm:"column:updated_by" json:"updated_by"`                           // 修改人
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"`                           // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at" json:"updated_at"`                           // 修改时间
	DeletedAt int64  `gorm:"column:deleted_at" json:"deleted_at"`                           // 删除时间
}
