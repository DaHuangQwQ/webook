package service

import (
	"github.com/gin-gonic/gin"
	"webook/internal/domain"
	"webook/internal/repository"
)

type RecruitmentService interface {
	Add(ctx *gin.Context, recruitment domain.Recruitment) error
}

type recruitmentService struct {
	repo repository.RecruitmentRepository
}

func NewRecruitmentService(repo repository.RecruitmentRepository) RecruitmentService {
	return &recruitmentService{
		repo: repo,
	}
}

func (svc *recruitmentService) Add(ctx *gin.Context, recruitment domain.Recruitment) error {
	return svc.repo.Create(ctx, recruitment)
}
