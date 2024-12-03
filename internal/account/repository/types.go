package repository

import (
	"context"
	"github.com/DaHuangQwQ/webook/internal/account/domain"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, c domain.Credit) error
}
