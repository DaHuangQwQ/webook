package api

import "webook/user/domain"

type RoleListReq struct {
	Meta     `path:"/role/list" method:"get" summary:"角色列表"`
	RoleName string `p:"roleName"`   //参数名称
	Status   string `p:"roleStatus"` //状态
	PageReq
}

type RoleListRes struct {
	ListRes
	List []domain.Role `json:"list"`
}
