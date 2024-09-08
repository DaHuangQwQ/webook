package domain

import "time"

type User struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`

	Phone    string `json:"mobile"`
	Password string `json:"password"`

	Nickname string `json:"userName"`
	Grade    int    `json:"grade"`

	Gender int    `json:"gender"`
	Avatar string `json:"avatar"`

	CTime time.Time `json:"ctime"`

	WechatInfo WechatInfo

	Birthday    int    `json:"birthday"      description:"生日"`
	UserStatus  uint   `json:"userStatus"    description:"用户状态;0:禁用,1:正常,2:未验证"`
	DeptId      uint64 `json:"deptId"        description:"部门id"`
	Remark      string `json:"remark"        description:"备注"`
	IsAdmin     int    `json:"isAdmin"       description:"是否后台管理员 1 是  0   否"`
	Address     string `json:"address"       description:"联系地址"`
	Describe    string `json:"describe"      description:"描述信息"`
	LastLoginIp string `json:"lastLoginIp"   description:"最后登录ip"`
}

type UserInfo struct {
	Nickname string `json:"username"`
	Grade    int    `json:"grade"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	// 0 未知， 1 男， 2 女
	Gender int `json:"gender"`
}
