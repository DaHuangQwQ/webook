package repository

import (
	"context"
	"github.com/DaHuangQwQ/webook/account/domain"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, c domain.Credit) error
}
