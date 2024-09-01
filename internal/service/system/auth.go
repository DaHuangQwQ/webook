package system

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/system"
)

type AuthService interface {
	Add(ctx context.Context, authRule domain.SysAuthRule) error
	Update(ctx context.Context, authRule domain.SysAuthRule) error
	Delete(ctx context.Context, ids []int64) error
	Get(ctx context.Context, id int64) (domain.SysAuthRule, error)
	List(ctx context.Context) ([]domain.SysAuthRule, error)
	GetIsMenuList(ctx context.Context) ([]*domain.SysAuthRule, error)
	GetIsButtonList(ctx context.Context) ([]*domain.SysAuthRule, error)
	GetMenuList(ctx context.Context) (list []*domain.SysAuthRule, err error)
	GetMenuRoles(ctx context.Context, id int64) (roleIds []uint, err error)
}

type authService struct {
	repo system.AuthRepository
}

func NewAuthService(repo system.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (svc *authService) GetMenuRoles(ctx context.Context, id int64) (roleIds []uint, err error) {
	return svc.repo.GetMenuRoles(ctx, id)
}

func (svc *authService) GetIsButtonList(ctx context.Context) ([]*domain.SysAuthRule, error) {
	list, err := svc.GetMenuList(ctx)
	if err != nil {
		return nil, err
	}
	var gList = make([]*domain.SysAuthRule, 0, len(list))
	for _, v := range list {
		if v.MenuType == 2 {
			gList = append(gList, v)
		}
	}
	return gList, nil
}

// GetMenuList 获取所有菜单
func (svc *authService) GetMenuList(ctx context.Context) ([]*domain.SysAuthRule, error) {
	listRes, err := svc.repo.FindAll(ctx)
	list := make([]*domain.SysAuthRule, len(listRes))
	for i, v := range listRes {
		list[i] = &v
	}
	return list, err
}

func (svc *authService) GetIsMenuList(ctx context.Context) ([]*domain.SysAuthRule, error) {
	list, err := svc.GetMenuList(ctx)
	if err != nil {
		return nil, err
	}
	var gList = make([]*domain.SysAuthRule, 0, len(list))
	for _, v := range list {
		if v.MenuType == 0 || v.MenuType == 1 {
			gList = append(gList, v)
		}
	}
	return gList, nil
}

func (svc *authService) Add(ctx context.Context, authRule domain.SysAuthRule) error {
	return svc.repo.Create(ctx, authRule)
}

func (svc *authService) Update(ctx context.Context, authRule domain.SysAuthRule) error {
	return svc.repo.Update(ctx, authRule)
}

func (svc *authService) Delete(ctx context.Context, ids []int64) error {
	return svc.repo.DeleteByIds(ctx, ids)
}

func (svc *authService) Get(ctx context.Context, id int64) (domain.SysAuthRule, error) {
	res, err := svc.repo.FindById(ctx, id)
	return res, err
}

func (svc *authService) List(ctx context.Context) ([]domain.SysAuthRule, error) {
	res, err := svc.repo.FindAll(ctx)
	return res, err
}
