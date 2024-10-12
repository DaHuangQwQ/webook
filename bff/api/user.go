package api

import (
	domain2 "webook/user/domain"
)

// UserMenusReq
// path:"/user/getUserMenus" tags:"用户管理" method:"get" summary:"获取用户菜单"
type UserMenusReq struct {
	//commonApi.Author
}

type UserMenusRes struct {
	MenuList    []domain2.UserMenus `json:"menuList"`
	Permissions []string            `json:"permissions"`
}

// UserSearchReq 用户搜索请求参数
// path:"/user/list" tags:"用户管理" method:"get" summary:"用户列表"
type UserSearchReq struct {
	Meta     `path:"/user/list" method:"get"`
	DeptId   string `json:"deptId"` //部门id
	Mobile   string `json:"mobile"`
	Status   string `json:"status"`
	KeyWords string `json:"keyWords"`
	PageReq
	//commonApi.Author
}

type UserSearchRes struct {
	UserList []SysUserRoleDeptRes `json:"userList"`
	ListRes
}

// UserGetParamsReq
// path:"/user/params" tags:"用户管理" method:"get" summary:"用户维护参数获取"
type UserGetParamsReq struct {
	Meta `path:"/user/params" method:"get"`
}

type UserGetParamsRes struct {
	RoleList []domain2.Role `json:"roleList"`
	//Posts    []*entity.SysPost `json:"posts"`
}

// SetUserReq 添加修改用户公用请求字段
type SetUserReq struct {
	Meta     `path:"/user/set" method:"post"`
	DeptId   uint64 `json:"deptId" v:"required#用户部门不能为空"` //所属部门
	Email    string `json:"email" v:"email#邮箱格式错误"`       //邮箱
	NickName string `json:"nickName" v:"required#用户昵称不能为空"`
	Mobile   string `json:"mobile" v:"required|phone#手机号不能为空|手机号格式错误"`
	//PostIds  []int64 `json:"postIds"`
	Remark  string  `json:"remark"`
	RoleIds []int64 `json:"roleIds"`
	Sex     int     `json:"sex"`
	Status  uint    `json:"status"`
	IsAdmin int     `json:"isAdmin"` // 是否后台管理员 1 是  0   否
}

// UserAddReq 添加用户参数
// path:"/user/add" tags:"用户管理" method:"post" summary:"添加用户"
type UserAddReq struct {
	Meta `path:"/user/add" method:"post"`
	*SetUserReq
	UserName string `json:"userName" v:"required#用户账号不能为空"`
	Password string `json:"password" v:"required|password#密码不能为空|密码以字母开头，只能包含字母、数字和下划线，长度在6~18之间"`
	UserSalt string
}

type UserAddRes struct {
}

// UserEditReq 修改用户参数
// path:"/user/edit" tags:"用户管理" method:"put" summary:"修改用户"
type UserEditReq struct {
	Meta `path:"/user/edit" method:"post"`
	*SetUserReq
	UserId int64 `json:"userId" v:"required#用户id不能为空"`
}

type UserEditRes struct {
}

// UserGetEditReq
// path:"/user/getEdit" tags:"用户管理" method:"get" summary:"获取用户信息"
type UserGetEditReq struct {
	Meta `path:"/user/getEdit" method:"get"`
	Id   uint64 `json:"id"`
}

type UserGetEditRes struct {
	//g.Meta         `mime:"application/json"`
	User           *domain2.User `json:"user"`
	CheckedRoleIds []uint        `json:"checkedRoleIds"`
	//CheckedPosts   []int64      `json:"checkedPosts"`
}

// UserResetPwdReq 重置用户密码状态参数
type UserResetPwdReq struct {
	//g.Meta   `path:"/user/resetPwd" tags:"用户管理" method:"put" summary:"重置用户密码"`
	Id       uint64 `json:"userId" v:"required#用户id不能为空"`
	Password string `json:"password" v:"required|password#密码不能为空|密码以字母开头，只能包含字母、数字和下划线，长度在6~18之间"`
}

type UserResetPwdRes struct {
}

// UserStatusReq 设置用户状态参数
// path:"/user/setStatus" tags:"用户管理" method:"put" summary:"设置用户状态"
type UserStatusReq struct {
	Meta       `path:"/user/setStatus" method:"post"`
	Id         uint64 `json:"userId" v:"required#用户id不能为空"`
	UserStatus uint   `json:"status" v:"required#用户状态不能为空"`
}

type UserStatusRes struct {
}

// UserDeleteReq
// path:"/user/delete" tags:"用户管理" method:"delete" summary:"删除用户"
type UserDeleteReq struct {
	Meta `path:"/user/delete" method:"delete"`
	Ids  []int `json:"ids"  v:"required#ids不能为空"`
}

type UserDeleteRes struct {
}

// UserGetByIdsReq
// path:"/user/getUsers" tags:"用户管理" method:"get" summary:"同时获取多个用户"
type UserGetByIdsReq struct {
	//commonApi.Author
	Ids []int `json:"ids" v:"required#ids不能为空"`
}

type UserGetByIdsRes struct {
	List []SysUserSimpleRes `json:"list"`
}

type SysUserSimpleRes struct {
	Id           uint64 `json:"id"`           //
	Avatar       string `json:"avatar"`       // 头像
	Sex          int    `json:"sex"`          // 性别
	UserName     string `json:"userName"`     // 用户名
	UserNickname string `json:"userNickname"` // 用户昵称
}

type SysUserRoleDeptRes struct {
	domain2.User
	Dept     domain2.SysDept      `json:"dept"`
	RoleInfo []SysUserRoleInfoRes `json:"roleInfo"`
}

type SysUserRoleInfoRes struct {
	RoleId uint   `json:"roleId"`
	Name   string `json:"name"`
}

type UserSignUpReq struct {
	Meta            `path:"/users/signup" method:"post"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UserLoginReq struct {
	Meta     `path:"/users/login" method:"post"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
