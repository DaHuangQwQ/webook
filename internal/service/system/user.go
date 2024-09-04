package system

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"webook/internal/api"
	"webook/internal/domain"
	"webook/internal/repository/system"
)

type UserService interface {
	GetAdminRules(ctx context.Context, userId int64) (userMenus []*domain.UserMenus, permissions []string, err error)
	GetAdminRole(ctx context.Context, userId int64, allRoleList []*domain.Role) ([]*domain.Role, error)
	GetAdminMenusByRoleIds(ctx context.Context, roleIds []uint) ([]*domain.UserMenus, error)
	GetPermissions(ctx context.Context, roleIds []uint) (permissions []string, err error)
	GetMenusTree(menus []*domain.UserMenus, pid uint) []*domain.UserMenus
	NotCheckAuthAdminIds(ctx context.Context)
	GetAllMenus(ctx context.Context) ([]*domain.UserMenus, error)
	GetAdminRoleIds(ctx context.Context, userId int64) (roleIds []uint, err error)
	GetUserSearch(ctx context.Context, req api.UserSearchReq) (api.UserSearchRes, error)
	Add(ctx *gin.Context, req api.UserAddReq) error
	Delete(ctx *gin.Context, ids []int) error
	GetParams(ctx *gin.Context) (api.UserGetParamsRes, error)
	GetEdit(ctx *gin.Context, id uint64) (api.UserGetEditRes, error)
}

type userService struct {
	roleSvc RoleService
	authSvc AuthService
	deptSvc DeptService
	repo    system.UserRepository
}

func NewSystemService(roleSvc RoleService, authSvc AuthService, repo system.UserRepository, deptSvc DeptService) UserService {
	return &userService{
		roleSvc: roleSvc,
		authSvc: authSvc,
		repo:    repo,
		deptSvc: deptSvc,
	}
}

func (svc *userService) GetEdit(ctx *gin.Context, id uint64) (api.UserGetEditRes, error) {
	//svc.repo.GetUserInfoById(ctx, id)
	return api.UserGetEditRes{}, nil
}

func (svc *userService) GetParams(ctx *gin.Context) (res api.UserGetParamsRes, err error) {
	roleList, err := svc.roleSvc.GetRoleList(ctx)
	res.RoleList = roleList
	return
}

func (svc *userService) Delete(ctx *gin.Context, ids []int) error {
	return svc.repo.DeleteByIds(ctx, ids)
}

func (svc *userService) Add(ctx *gin.Context, req api.UserAddReq) error {
	return svc.repo.Add(ctx, req)
}

func (svc *userService) GetUserSearch(ctx context.Context, req api.UserSearchReq) (res api.UserSearchRes, err error) {
	total, userList, err := svc.repo.List(ctx, req)
	if err != nil || total == 0 {
		return
	}
	res.Total = total
	allRoles, err := svc.roleSvc.GetRoleList(ctx)
	if err != nil {
		return
	}
	allRolesTemp := make([]*domain.Role, len(allRoles))
	for k, r := range allRoles {
		allRolesTemp[k] = &r
	}

	allDepts, err := svc.deptSvc.GetDeptList(ctx)
	if err != nil {
		return
	}
	users := make([]api.SysUserRoleDeptRes, len(userList))

	for k, u := range userList {
		var dept domain.SysDept
		users[k] = api.SysUserRoleDeptRes{
			User: u,
		}
		for _, d := range allDepts {
			if u.DeptId == d.DeptId {
				dept = d
			}
		}
		users[k].Dept = dept

		//roles, err := svc.GetAdminRole(ctx, u.Id, allRolesTemp)

		//if err != nil {
		//	return res, err
		//}
		for _, r := range allRolesTemp {
			users[k].RoleInfo = append(users[k].RoleInfo, api.SysUserRoleInfoRes{RoleId: uint(r.Id), Name: r.Name})
		}
	}
	res.UserList = users
	return
}

func (svc *userService) GetAdminRoleIds(ctx context.Context, userId int64) (roleIds []uint, err error) {
	return svc.repo.GetAdminRoleIds(ctx, userId)
}

func (svc *userService) setMenuData(menu *domain.UserMenu, entity *domain.SysAuthRule) *domain.UserMenu {
	menu = &domain.UserMenu{
		Id:   entity.Id,
		Pid:  entity.Pid,
		Name: camelLower(strings.Replace(entity.Name, "/", "_", -1)),
		//Name:      gstr.CaseCamelLower(gstr.Replace(entity.Name, "/", "_")),
		Component: entity.Component,
		Path:      entity.Path,
		MenuMeta: &domain.MenuMeta{
			Icon:        entity.Icon,
			Title:       entity.Title,
			IsLink:      "",
			IsHide:      entity.IsHide == 1,
			IsKeepAlive: entity.IsCached == 1,
			IsAffix:     entity.IsAffix == 1,
			IsIframe:    entity.IsIframe == 1,
		},
	}
	if menu.MenuMeta.IsIframe || entity.IsLink == 1 {
		menu.MenuMeta.IsLink = entity.LinkUrl
	}
	return menu
}

func (svc *userService) GetAllMenus(ctx context.Context) (menus []*domain.UserMenus, err error) {
	//获取所有开启的菜单
	var allMenus []*domain.SysAuthRule
	allMenus, err = svc.authSvc.GetIsMenuList(ctx)
	if err != nil {
		return
	}
	menus = make([]*domain.UserMenus, len(allMenus))
	for k, v := range allMenus {
		var menu *domain.UserMenu
		menu = svc.setMenuData(menu, v)
		menus[k] = &domain.UserMenus{UserMenu: menu}
	}
	menus = svc.GetMenusTree(menus, 0)
	return
}

// NotCheckAuthAdminIds super admin
func (svc *userService) NotCheckAuthAdminIds(ctx context.Context) {
	return
}

func (svc *userService) GetAdminRules(ctx context.Context, userId int64) (menuList []*domain.UserMenus, permissions []string, err error) {
	//是否超管
	isSuperAdmin := true
	//获取无需验证权限的用户id
	//svc.NotCheckAuthAdminIds()

	//获取用户菜单数据
	roleListRes, err := svc.roleSvc.GetRoleList(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("GetRoleList失败, %w", err)
	}
	roleList := make([]*domain.Role, len(roleListRes))
	for i, role := range roleListRes {
		roleList[i] = &role
	}
	adminRoleList, err := svc.GetAdminRole(ctx, userId, roleList)
	if err != nil {
		return nil, nil, fmt.Errorf("GetAdminRole失败, %w", err)
	}
	name := make([]string, len(adminRoleList))
	roleIds := make([]uint, len(adminRoleList))
	for k, v := range adminRoleList {
		name[k] = v.Name
		roleIds[k] = uint(v.Id)
	}
	//获取菜单信息
	if isSuperAdmin {
		//超管获取所有菜单
		permissions = []string{"*/*/*"}
		menuList, err = svc.GetAllMenus(ctx)

	} else {
		menuList, err = svc.GetAdminMenusByRoleIds(ctx, roleIds)
		if err != nil {
			return nil, nil, fmt.Errorf("GetAdminMenusByRoleIds失败, %w", err)
		}
		permissions, err = svc.GetPermissions(ctx, roleIds)
		if err != nil {
			return nil, nil, fmt.Errorf("GetPermissions失败, %w", err)
		}
	}
	return
}

func (svc *userService) GetAdminRole(ctx context.Context, userId int64, allRoleList []*domain.Role) (roles []*domain.Role, err error) {
	roleIds, err := svc.GetAdminRoleIds(ctx, userId)
	if err != nil {
		err = fmt.Errorf("GetAdminRoleIds失败 %w", err)
		return
	}
	roles = make([]*domain.Role, 0, len(allRoleList))
	for _, v := range allRoleList {
		for _, id := range roleIds {
			if int64(id) == v.Id {
				roles = append(roles, v)
			}
		}
		if len(roles) == len(roleIds) {
			break
		}
	}
	return
}

func (svc *userService) GetAdminMenusByRoleIds(ctx context.Context, roleIds []uint) (menus []*domain.UserMenus, err error) {
	//获取角色对应的菜单id
	menuIds := svc.repo.GetMenuIds(ctx, roleIds)
	//获取所有开启的菜单
	allMenus, err := svc.authSvc.GetIsMenuList(ctx)
	for _, v := range allMenus {
		if _, ok := menuIds[int64(v.Id)]; (v.Condition == "nocheck") || ok {
			var roleMenu *domain.UserMenu
			roleMenu = svc.setMenuData(roleMenu, v)
			menus = append(menus, &domain.UserMenus{UserMenu: roleMenu})
		}
	}
	menus = svc.GetMenusTree(menus, 0)
	return
}

func (svc *userService) GetPermissions(ctx context.Context, roleIds []uint) (userButtons []string, err error) {
	menuIds := svc.repo.GetMenuIds(ctx, roleIds)
	//获取所有开启的按钮
	allButtons, err := svc.authSvc.GetIsButtonList(ctx)
	if err != nil {
		return nil, err
	}
	userButtons = make([]string, 0, len(allButtons))
	for _, button := range allButtons {
		if _, ok := menuIds[int64(button.Id)]; (button.Condition == "nocheck") || ok {
			userButtons = append(userButtons, button.Name)
		}
	}
	return
}

func (svc *userService) GetMenusTree(menus []*domain.UserMenus, pid uint) []*domain.UserMenus {
	returnList := make([]*domain.UserMenus, 0, len(menus))
	for _, menu := range menus {
		if menu.Pid == pid {
			menu.Children = svc.GetMenusTree(menus, menu.Id)
			returnList = append(returnList, menu)
		}
	}
	return returnList
}

// camelLower 将字符串从任何形式（假设用 "_" 分隔单词）转换为驼峰命名（但首字母小写）
func camelLower(s string) string {
	// 首先，将所有的 "_" 替换为空格，方便之后分割
	s = strings.Replace(s, "_", " ", -1)

	// 去除字符串两头的空格
	s = strings.TrimSpace(s)

	// 分割字符串为单词数组
	words := strings.Fields(s)

	// 如果没有单词，直接返回空字符串
	if len(words) == 0 {
		return ""
	}

	// 初始化结果字符串，并将第一个单词的首字母转换为小写
	camel := strings.ToLower(string(words[0][0])) + words[0][1:]

	// 遍历剩余的单词，将每个单词的首字母转换为大写并附加到结果字符串
	for _, word := range words[1:] {
		camel += strings.Title(word)
	}

	return camel
}
