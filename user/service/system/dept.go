package system

import (
	"context"
	"webook/user/domain"
	"webook/user/repository/system"
)

type DeptService interface {
	List(ctx context.Context, deptName string, status uint, pageNum, pageSize int) ([]domain.SysDept, error)
	Add(ctx context.Context, dept domain.SysDept) error
	Edit(ctx context.Context, dept domain.SysDept) error
	Delete(ctx context.Context, deptId int64) error
	TreeSelect(ctx context.Context) (domain.DeptTreeSelectRes, error)
	GetDeptList(ctx context.Context) ([]domain.SysDept, error)
	GetListTree(pid uint64, list []domain.SysDept) []*domain.SysDeptTreeRes
}

type deptService struct {
	repo system.DeptRepository
}

func NewDeptService(repo system.DeptRepository) DeptService {
	return &deptService{
		repo: repo,
	}
}

func (svc *deptService) GetDeptList(ctx context.Context) ([]domain.SysDept, error) {
	return svc.repo.GetAllList(ctx)
}

func (svc *deptService) List(ctx context.Context, deptName string, status uint, pageNum, pageSize int) ([]domain.SysDept, error) {
	return svc.repo.GetList(ctx, deptName, status, pageNum, pageSize)
}

func (svc *deptService) Add(ctx context.Context, dept domain.SysDept) error {
	//TODO implement me
	panic("implement me")
}

func (svc *deptService) Edit(ctx context.Context, dept domain.SysDept) error {
	//TODO implement me
	panic("implement me")
}

func (svc *deptService) Delete(ctx context.Context, deptId int64) error {
	//TODO implement me
	panic("implement me")
}

func (svc *deptService) TreeSelect(ctx context.Context) (domain.DeptTreeSelectRes, error) {
	var res domain.DeptTreeSelectRes
	deptList, err := svc.GetDeptList(ctx)
	if err != nil {
		return domain.DeptTreeSelectRes{}, nil
	}
	res.Deps = svc.GetListTree(0, deptList)
	return res, err
}

func (svc *deptService) GetListTree(pid uint64, list []domain.SysDept) []*domain.SysDeptTreeRes {
	deptTree := make([]*domain.SysDeptTreeRes, 0, len(list))
	for _, v := range list {
		if v.ParentId == pid {
			t := &domain.SysDeptTreeRes{
				SysDept: &v,
			}
			child := svc.GetListTree(v.DeptId, list)
			if len(child) > 0 {
				t.Children = child
			}
			deptTree = append(deptTree, t)
		}
	}
	return deptTree
}
