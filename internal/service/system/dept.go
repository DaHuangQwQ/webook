package system

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/system"
)

type DeptService interface {
	List(ctx context.Context, deptName string, status uint, pageNum, pageSize int) ([]domain.SysDept, error)
	Add(ctx context.Context, dept domain.SysDept) error
	Edit(ctx context.Context, dept domain.SysDept) error
	Delete(ctx context.Context, deptId int64) error
	TreeSelect(ctx context.Context) []domain.SysDeptTreeRes
}

type deptService struct {
	repo system.DeptRepository
}

func NewDeptService(repo system.DeptRepository) DeptService {
	return &deptService{
		repo: repo,
	}
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

func (svc *deptService) TreeSelect(ctx context.Context) []domain.SysDeptTreeRes {
	//TODO implement me
	panic("implement me")
}
