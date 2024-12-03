package service

import (
	"context"
	"github.com/DaHuangQwQ/webook/internal/account/domain"
	"github.com/DaHuangQwQ/webook/internal/account/repository"
)

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (a *accountService) Credit(ctx context.Context, cr domain.Credit) error {
	return a.repo.AddCredit(ctx, cr)
}
