package system

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"strconv"
)

type UserRepository interface {
	GetMenuIds(ctx context.Context, roleIds []uint) map[int64]int64
	GetAdminRoleIds(ctx context.Context, userId int64) (roleIds []uint, err error)
}

type CachedUserRepository struct {
	casbin           casbin.IEnforcer
	casBinUserPrefix string
}

func NewCachedUserRepository(casbin casbin.IEnforcer) UserRepository {
	return &CachedUserRepository{
		casbin:           casbin,
		casBinUserPrefix: "u_",
	}
}

func (repo *CachedUserRepository) GetAdminRoleIds(ctx context.Context, userId int64) (roleIds []uint, err error) {
	groupPolicy, err := repo.casbin.GetFilteredGroupingPolicy(0, fmt.Sprintf("%s%d", repo.casBinUserPrefix, userId))
	if len(groupPolicy) > 0 {
		roleIds = make([]uint, len(groupPolicy))
		//得到角色id的切片

		for k, v := range groupPolicy {
			num, _ := strconv.ParseUint(v[1], 10, 64)
			roleIds[k] = uint(num)
		}
	}
	return
}

func (repo *CachedUserRepository) GetMenuIds(ctx context.Context, roleIds []uint) map[int64]int64 {
	menuIds := map[int64]int64{}
	for _, roleId := range roleIds {
		//查询当前权限
		gp, _ := repo.casbin.GetFilteredPolicy(0, strconv.Itoa(int(roleId)))
		for _, p := range gp {
			mid, _ := strconv.ParseInt(p[1], 10, 64)
			menuIds[mid] = mid
		}
	}
	return menuIds
}
