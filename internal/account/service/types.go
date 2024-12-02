package service

import (
	"context"
	"github.com/DaHuangQwQ/webook/account/domain"
)

type AccountService interface {
	Credit(ctx context.Context, cr domain.Credit) error
}
