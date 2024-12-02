package repository

import (
	"github.com/DaHuangQwQ/webook/internal_temp/domain"
	"github.com/DaHuangQwQ/webook/internal_temp/repository/dao"
	"github.com/gin-gonic/gin"
)

type RecruitmentRepository interface {
	Create(ctx *gin.Context, recruitment domain.Recruitment) error
}

type CachedRecruitmentRepository struct {
	dao dao.RecruitmentDao
}

func NewCachedRecruitmentRepository(dao dao.RecruitmentDao) RecruitmentRepository {
	return &CachedRecruitmentRepository{
		dao: dao,
	}
}

func (repo *CachedRecruitmentRepository) Create(ctx *gin.Context, recruitment domain.Recruitment) error {
	return repo.dao.Create(ctx, repo.toEntity(recruitment))
}

func (repo *CachedRecruitmentRepository) toEntity(recruitment domain.Recruitment) dao.Recruitment {
	return dao.Recruitment{
		Id:          recruitment.Id,
		Name:        recruitment.Name,
		StudentID:   recruitment.StudentID,
		Major:       recruitment.Major,
		Situation:   recruitment.Situation,
		Expectation: recruitment.Expectation,
		Selfie:      recruitment.Selfie,
		ErrorNum:    recruitment.ErrorNum,
	}
}
