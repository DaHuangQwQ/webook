package system

import (
	"context"
	"github.com/casbin/casbin/v2"
	"strconv"
	"webook/user/domain"
	"webook/user/repository/dao/system"
)

type RoleRepository interface {
	GetRoleList(ctx context.Context) ([]domain.Role, error)
	DeleteRoleRule(ctx context.Context, id int64) error
	FindById(ctx context.Context, id int64) (domain.Role, error)
	FindAll(ctx context.Context) ([]domain.Role, error)
	Save(ctx context.Context, role domain.Role) error
	AddRoleRule(ctx context.Context, roleId int64, ruleIds []int64) error
	AddRole(ctx context.Context, role domain.Role) error
	GetRole(ctx context.Context, id int64) (domain.Role, error)
	EditRole(ctx context.Context, role domain.Role) error
	DeleteByIds(ctx context.Context, ids []int64) error
	GetRoleListSearch(ctx context.Context, role domain.Role, pageNum int, pageSize int) (int64, []domain.Role, error)
	GetFilteredNamedPolicy(ctx context.Context, id int64) (gpSlice []int, err error)
}

type CachedRoleRepository struct {
	casbin casbin.IEnforcer
	dao    system.RoleDao
}

func NewCachedRoleRepository(casbin casbin.IEnforcer, dao system.RoleDao) RoleRepository {
	return &CachedRoleRepository{
		casbin: casbin,
		dao:    dao,
	}
}

func (c *CachedRoleRepository) GetFilteredNamedPolicy(ctx context.Context, id int64) (gpSlice []int, err error) {
	gp, err := c.casbin.GetFilteredNamedPolicy("p", 0, strconv.FormatInt(id, 10))
	gpSlice = make([]int, len(gp))
	for k, v := range gp {
		i64, _ := strconv.ParseInt(v[1], 10, 64)
		gpSlice[k] = int(i64)
	}
	return
}

func (c *CachedRoleRepository) GetRoleListSearch(ctx context.Context, role domain.Role, pageNum int, pageSize int) (int64, []domain.Role, error) {
	total, res, err := c.dao.GetRoleListSearch(ctx, c.toRole(role), pageNum, pageSize)
	Role := make([]domain.Role, len(res))
	for i, v := range res {
		Role[i] = c.toDomain(v)
	}
	return total, Role, err
}

func (c *CachedRoleRepository) DeleteByIds(ctx context.Context, ids []int64) error {
	return c.dao.DeleteByIds(ctx, ids)
}

func (c *CachedRoleRepository) EditRole(ctx context.Context, role domain.Role) error {
	err := c.dao.UpdateById(ctx, c.toRole(role))
	if err != nil {
		return err
	}
	err = c.DeleteRoleRule(ctx, role.Id)
	if err != nil {
		return err
	}
	return c.AddRoleRule(ctx, role.Id, role.MenuIds)
}

func (c *CachedRoleRepository) GetRole(ctx context.Context, id int64) (domain.Role, error) {
	role, err := c.dao.FindById(ctx, id)
	return c.toDomain(role), err
}

func (c *CachedRoleRepository) AddRole(ctx context.Context, role domain.Role) error {
	id, err := c.dao.CreateAndGetId(ctx, c.toRole(role))
	if err != nil {
		return err
	}
	return c.AddRoleRule(ctx, id, role.MenuIds)
}

func (c *CachedRoleRepository) AddRoleRule(ctx context.Context, roleId int64, ruleIds []int64) error {
	for _, role_id := range ruleIds {
		_, err := c.casbin.AddPolicy(strconv.FormatInt(roleId, 10), strconv.FormatInt(role_id, 10), "All")
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CachedRoleRepository) GetRoleList(ctx context.Context) ([]domain.Role, error) {
	role, err := c.dao.FindList(ctx)
	res := make([]domain.Role, len(role))
	for i, role := range role {
		res[i] = c.toDomain(role)
	}
	return res, err
}

func (c *CachedRoleRepository) DeleteRoleRule(ctx context.Context, id int64) error {
	_, err := c.casbin.RemoveFilteredPolicy(0, strconv.FormatInt(id, 10))
	return err
}

func (c *CachedRoleRepository) FindById(ctx context.Context, id int64) (domain.Role, error) {
	res, err := c.dao.FindById(ctx, id)
	return c.toDomain(res), err
}

func (c *CachedRoleRepository) FindAll(ctx context.Context) ([]domain.Role, error) {
	return []domain.Role{}, nil
}

func (c *CachedRoleRepository) Save(ctx context.Context, role domain.Role) error {
	res, err := c.casbin.AddPolicy(strconv.FormatInt(role.Id, 10), strconv.FormatInt(role.Ids, 10), "All")
	//res, err := c.casbin.AddPolicy("1", "2", "All")
	if err != nil {
		return err
	}
	println(res)
	return nil
}

func (c *CachedRoleRepository) toDomain(role system.SysRole) domain.Role {
	return domain.Role{
		Id:        role.ID,
		Name:      role.Name,
		ListOrder: role.ListOrder,
		Status:    role.Status,
		Remark:    role.Remark,
		DataScope: role.DataScope,
	}
}

func (c *CachedRoleRepository) toRole(role domain.Role) system.SysRole {
	return system.SysRole{
		ID:        role.Id,
		Name:      role.Name,
		ListOrder: role.ListOrder,
		Status:    role.Status,
		Remark:    role.Remark,
		DataScope: role.DataScope,
	}

}
