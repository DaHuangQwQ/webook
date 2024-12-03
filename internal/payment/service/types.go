package service

import (
	"context"
	"github.com/DaHuangQwQ/webook/internal/payment/domain"
)

type PaymentService interface {
	Prepay(ctx context.Context, payment domain.Payment) (string, error)
}
