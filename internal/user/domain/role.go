package domain

type Role struct {
	Id        int64 `json:"id"        description:""`
	Status    uint8 `json:"status"    description:"状态;0:禁用;1:正常"`
	Ids       int64
	ListOrder uint   `json:"listOrder" description:"排序"`
	Name      string `json:"name"      description:"角色名称"`
	Remark    string `json:"remark"    description:"备注"`
	DataScope uint8  `json:"dataScope" description:"数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）"`
	MenuIds   []int64
}

type SysDept struct {
	DeptId    uint64 `json:"deptId"    description:"部门id"`
	ParentId  uint64 `json:"parentId"  description:"父部门id"`
	Ancestors string `json:"ancestors" description:"祖级列表"`
	DeptName  string `json:"deptName"  description:"部门名称"`
	OrderNum  int    `json:"orderNum"  description:"显示顺序"`
	Leader    string `json:"leader"    description:"负责人"`
	Phone     string `json:"phone"     description:"联系电话"`
	Email     string `json:"email"     description:"邮箱"`
	Status    uint   `json:"status"    description:"部门状态（0正常 1停用）"`
}

type SysAuthRule struct {
	Id         uint   `json:"id"         description:""`
	Pid        uint   `json:"pid"        description:"父ID"`
	Name       string `json:"name"       description:"规则名称"`
	Title      string `json:"title"      description:"规则名称"`
	Icon       string `json:"icon"       description:"图标"`
	Condition  string `json:"condition"  description:"条件"`
	Remark     string `json:"remark"     description:"备注"`
	MenuType   uint8  `json:"menuType"   description:"类型 0目录 1菜单 2按钮"`
	Weigh      int    `json:"weigh"      description:"权重"`
	IsHide     uint8  `json:"isHide"     description:"显示状态"`
	Path       string `json:"path"       description:"路由地址"`
	Component  string `json:"component"  description:"组件路径"`
	IsLink     uint8  `json:"isLink"     description:"是否外链 1是 0否"`
	ModuleType string `json:"moduleType" description:"所属模块"`
	ModelId    uint   `json:"modelId"    description:"模型ID"`
	IsIframe   uint8  `json:"isIframe"   description:"是否内嵌iframe"`
	IsCached   uint8  `json:"isCached"   description:"是否缓存"`
	Redirect   string `json:"redirect"   description:"路由重定向地址"`
	IsAffix    uint8  `json:"isAffix"    description:"是否固定"`
	LinkUrl    string `json:"linkUrl"    description:"链接地址"`

	Roles []uint `json:"roles"` // 角色ids

}

type UserMenu struct {
	Id        uint   `json:"id"`
	Pid       uint   `json:"pid"`
	Name      string `json:"name"`
	Component string `json:"component"`
	Path      string `json:"path"`
	*MenuMeta `json:"meta"`
}

type UserMenus struct {
	*UserMenu `json:""`
	Children  []*UserMenus `json:"children"`
}

type MenuMeta struct {
	Icon        string `json:"icon"`
	Title       string `json:"title"`
	IsLink      string `json:"isLink"`
	IsHide      bool   `json:"isHide"`
	IsKeepAlive bool   `json:"isKeepAlive"`
	IsAffix     bool   `json:"isAffix"`
	IsIframe    bool   `json:"isIframe"`
}

type SysDeptTreeRes struct {
	*SysDept
	Children []*SysDeptTreeRes `json:"children"`
}

type DeptTreeSelectRes struct {
	Deps []*SysDeptTreeRes `json:"deps"`
}
