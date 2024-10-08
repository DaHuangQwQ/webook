package system

import (
	"context"
	"webook/user/domain"
	"webook/user/repository/dao"
)

type DeptRepository interface {
	GetList(ctx context.Context, deptName string, status uint, pageNum, pageSize int) ([]domain.SysDept, error)
	GetAllList(ctx context.Context) ([]domain.SysDept, error)
}

type CachedDeptRepository struct {
	dao dao.DeptDao
}

func NewCachedDeptRepository(dao dao.DeptDao) DeptRepository {
	return &CachedDeptRepository{
		dao: dao,
	}
}

func (repo *CachedDeptRepository) GetAllList(ctx context.Context) ([]domain.SysDept, error) {
	res, err := repo.dao.FindAll(ctx)
	depts := make([]domain.SysDept, len(res))
	for i, dept := range res {
		depts[i] = repo.toDomain(dept)
	}
	return depts, err
}

func (repo *CachedDeptRepository) GetList(ctx context.Context, deptName string, status uint, pageNum, pageSize int) ([]domain.SysDept, error) {
	res, err := repo.dao.GetList(ctx, deptName, status, pageNum, pageSize)
	depts := make([]domain.SysDept, len(res))
	for i, item := range res {
		depts[i] = repo.toDomain(item)
	}
	return depts, err
}

func (repo *CachedDeptRepository) toDomain(dept dao.SysDept) domain.SysDept {
	return domain.SysDept{
		DeptId:    uint64(dept.DeptID),
		DeptName:  dept.DeptName,
		ParentId:  uint64(dept.ParentID),
		Status:    uint(dept.Status),
		Ancestors: dept.Ancestors,
		OrderNum:  dept.OrderNum,
		Leader:    dept.Leader,
		Phone:     dept.Phone,
		Email:     dept.Email,
	}
}

func (repo *CachedDeptRepository) roEntity(dept domain.SysDept) dao.SysDept {
	return dao.SysDept{
		DeptID:    int64(dept.DeptId),
		DeptName:  dept.DeptName,
		ParentID:  int64(dept.ParentId),
		Status:    uint8(dept.Status),
		Ancestors: dept.Ancestors,
		OrderNum:  dept.OrderNum,
		Leader:    dept.Leader,
		Phone:     dept.Phone,
		Email:     dept.Email,
	}
}
