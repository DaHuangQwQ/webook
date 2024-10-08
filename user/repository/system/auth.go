package system

import (
	"context"
	"github.com/casbin/casbin/v2"
	"strconv"
	"webook/user/domain"
	"webook/user/repository/dao"
)

type AuthRepository interface {
	FindById(ctx context.Context, id int64) (domain.SysAuthRule, error)
	FindAll(ctx context.Context) ([]domain.SysAuthRule, error)
	Create(ctx context.Context, authRule domain.SysAuthRule) error
	DeleteByIds(ctx context.Context, ids []int64) error
	Update(ctx context.Context, authRule domain.SysAuthRule) error
	GetMenuRoles(ctx context.Context, id int64) (roleIds []uint, err error)
}

type CachedAuthRepository struct {
	dao    dao.AuthDao
	casbin casbin.IEnforcer
}

func NewCachedAuthRepository(dao dao.AuthDao, casbin casbin.IEnforcer) AuthRepository {
	return &CachedAuthRepository{
		dao:    dao,
		casbin: casbin,
	}
}

func (repo *CachedAuthRepository) GetMenuRoles(ctx context.Context, id int64) (roleIds []uint, err error) {
	policies, err := repo.casbin.GetFilteredNamedPolicy("p", 1, strconv.FormatInt(id, 10))
	for _, policy := range policies {
		parseUint, _ := strconv.ParseUint(policy[0], 10, 64)
		roleIds = append(roleIds, uint(parseUint))
	}
	return
}

func (repo *CachedAuthRepository) FindById(ctx context.Context, id int64) (domain.SysAuthRule, error) {
	res, err := repo.dao.Find(ctx, id)
	return repo.toDomain(res), err
}

func (repo *CachedAuthRepository) FindAll(ctx context.Context) ([]domain.SysAuthRule, error) {
	res, err := repo.dao.FindAll(ctx)
	authRule := make([]domain.SysAuthRule, len(res))
	for i, r := range res {
		authRule[i] = repo.toDomain(r)
	}
	return authRule, err
}

func (repo *CachedAuthRepository) Create(ctx context.Context, authRule domain.SysAuthRule) error {
	return repo.dao.Insert(ctx, repo.toEntity(authRule))
}

func (repo *CachedAuthRepository) DeleteByIds(ctx context.Context, ids []int64) error {
	return repo.dao.DeleteByIds(ctx, ids)
}

func (repo *CachedAuthRepository) Update(ctx context.Context, authRule domain.SysAuthRule) error {
	return repo.dao.Update(ctx, repo.toEntity(authRule))
}

func (repo *CachedAuthRepository) toEntity(authRule domain.SysAuthRule) dao.SysAuthRule {
	return dao.SysAuthRule{
		ID:         authRule.Id,
		PID:        authRule.Pid,
		Name:       authRule.Name,
		Title:      authRule.Title,
		Icon:       authRule.Icon,
		Condition:  authRule.Condition,
		Remark:     authRule.Remark,
		MenuType:   authRule.MenuType,
		Weigh:      authRule.Weigh,
		IsHide:     authRule.IsHide,
		Path:       authRule.Path,
		Component:  authRule.Component,
		IsLink:     authRule.IsLink,
		ModuleType: authRule.ModuleType,
		ModelID:    authRule.ModelId,
		IsIframe:   authRule.IsIframe,
		IsCached:   authRule.IsCached,
		Redirect:   authRule.Redirect,
		IsAffix:    authRule.IsAffix,
		LinkURL:    authRule.LinkUrl,
	}
}

func (repo *CachedAuthRepository) toDomain(authRule dao.SysAuthRule) domain.SysAuthRule {
	return domain.SysAuthRule{
		Id:         authRule.ID,
		Pid:        authRule.PID,
		Name:       authRule.Name,
		Title:      authRule.Title,
		Icon:       authRule.Icon,
		Condition:  authRule.Condition,
		Remark:     authRule.Remark,
		MenuType:   authRule.MenuType,
		Weigh:      authRule.Weigh,
		IsHide:     authRule.IsHide,
		Path:       authRule.Path,
		Component:  authRule.Component,
		IsLink:     authRule.IsLink,
		ModuleType: authRule.ModuleType,
		ModelId:    authRule.ModelID,
		IsIframe:   authRule.IsIframe,
		IsCached:   authRule.IsCached,
		Redirect:   authRule.Redirect,
		IsAffix:    authRule.IsAffix,
		LinkUrl:    authRule.LinkURL,
	}
}
