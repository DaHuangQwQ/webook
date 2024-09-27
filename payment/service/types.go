package service

import (
	"context"
	"webook/payment/domain"
)

type PaymentService interface {
	Prepay(ctx context.Context, payment domain.Payment) (string, error)
}
