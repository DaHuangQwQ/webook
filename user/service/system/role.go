package system

import (
	"context"
	"webook/user/api"
	"webook/user/domain"
	system2 "webook/user/repository/system"
)

type RoleService interface {
	GetRoleListSearch(ctx context.Context, role domain.Role, page api.PageReq) (int64, []domain.Role, error)
	GetRoleList(ctx context.Context) ([]domain.Role, error)
	AddRoleRule(ctx context.Context, ruleIds []int64, roleId int64) (err error)
	DelRoleRule(ctx context.Context, roleId int64) error
	AddRole(ctx context.Context, role domain.Role) error
	Get(ctx context.Context, id int64) (domain.Role, error)
	GetFilteredNamedPolicy(ctx context.Context, id int64) (gpSlice []int, err error)
	EditRole(ctx context.Context, role domain.Role) error
	DeleteByIds(ctx context.Context, ids []int64) error
	GetParams(ctx context.Context) ([]domain.SysAuthRule, error)
}

type role struct {
	repo     system2.RoleRepository
	authRepo system2.AuthRepository
}

func NewRoleService(repo system2.RoleRepository, authRepo system2.AuthRepository) RoleService {
	return &role{
		repo:     repo,
		authRepo: authRepo,
	}
}

func (r *role) GetParams(ctx context.Context) ([]domain.SysAuthRule, error) {
	return r.authRepo.FindAll(ctx)
}

func (r *role) GetRoleListSearch(ctx context.Context, role domain.Role, page api.PageReq) (int64, []domain.Role, error) {
	total, roles, err := r.repo.GetRoleListSearch(ctx, role, page.PageNum, page.PageSize)
	if err != nil {
		return 0, nil, err
	}
	return total, roles, err
}

func (r *role) GetRoleList(ctx context.Context) ([]domain.Role, error) {
	return r.repo.GetRoleList(ctx)
}

func (r *role) AddRoleRule(ctx context.Context, ruleIds []int64, roleId int64) (err error) {
	return r.repo.AddRoleRule(ctx, roleId, ruleIds)
}

func (r *role) DelRoleRule(ctx context.Context, roleId int64) error {
	return r.repo.DeleteRoleRule(ctx, roleId)
}

func (r *role) AddRole(ctx context.Context, role domain.Role) error {
	return r.repo.AddRole(ctx, role)
}

func (r *role) Get(ctx context.Context, id int64) (domain.Role, error) {
	return r.repo.FindById(ctx, id)
}

func (r *role) GetFilteredNamedPolicy(ctx context.Context, id int64) (gpSlice []int, err error) {
	return r.repo.GetFilteredNamedPolicy(ctx, id)
}

func (r *role) EditRole(ctx context.Context, role domain.Role) error {
	return r.repo.EditRole(ctx, role)
}

func (r *role) DeleteByIds(ctx context.Context, ids []int64) error {
	return r.repo.DeleteByIds(ctx, ids)
}
