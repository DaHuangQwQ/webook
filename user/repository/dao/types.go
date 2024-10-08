package dao

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"time"
	"webook/bff/api"
)

type UserDao interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
	Update(ctx context.Context, user User) error
	FindAll(ctx context.Context, req api.UserSearchReq) (total int, userList []User, err error)
	InsertAndGetId(ctx context.Context, user User) (int64, error)
	DeleteByIds(ctx context.Context, ids []int) error
	UpdateStatus(ctx context.Context, user User) error
}

type AuthDao interface {
	Insert(ctx context.Context, authRule SysAuthRule) error
	Update(ctx context.Context, authRule SysAuthRule) error
	DeleteByIds(ctx context.Context, ids []int64) error
	Find(ctx context.Context, id int64) (SysAuthRule, error)
	FindAll(ctx context.Context) ([]SysAuthRule, error)
}

type DeptDao interface {
	Insert(ctx context.Context, dept SysDept) error
	Update(ctx context.Context, dept SysDept) error
	Delete(ctx context.Context, dept SysDept) error
	Find(ctx context.Context, dept SysDept) (SysDept, error)
	FindAll(ctx context.Context) ([]SysDept, error)
	GetList(ctx context.Context, deptName string, status uint, pageNum, pageSize int) ([]SysDept, error)
}

type RoleDao interface {
	FindList(ctx context.Context) ([]SysRole, error)
	UpdateById(ctx context.Context, role SysRole) error
	Create(ctx context.Context, role SysRole) error
	FindById(ctx context.Context, roleId int64) (SysRole, error)
	DeleteByIds(ctx context.Context, roleIds []int64) error
	CreateAndGetId(ctx context.Context, role SysRole) (id int64, err error)
	GetRoleListSearch(ctx context.Context, role SysRole, pageNum int, pageSize int) (int64, []SysRole, error)
}

type User struct {
	Id int64 `json:"id" gorm:"primary_key;autoIncrement"`

	WechatOpenId  sql.NullString `json:"wechat_open_id" gorm:"type:varchar(255);unique;"`
	WechatUnionId sql.NullString `json:"wechat_union_id" gorm:"type:varchar(255);unique;"`

	Email      sql.NullString `json:"email" gorm:"type:varchar(100);unique;sql:null"`
	Phone      sql.NullString `json:"phone" gorm:"type:varchar(100);unique;sql:null"`
	Password   string         `json:"password" gorm:"type:varchar(100)"`
	Nickname   string         `json:"nickname" gorm:"type:varchar(100)"`
	Birthday   time.Time      `json:"birthday" gorm:"type:datetime"`
	UserStatus uint8          `json:"user_status" gorm:"type:tinyint unsigned;not null;default:1;comment:'用户状态;0:禁用,1:正常,2:未验证'"`

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

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&SysRole{}, &SysDept{}, &User{}, &SysAuthRule{})
}
