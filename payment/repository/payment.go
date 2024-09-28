package repository

import (
	"context"
	"database/sql"
	"time"
	"webook/payment/domain"
	"webook/payment/repository/dao"
)

type paymentRepository struct {
	dao dao.PaymentDAO
}

func NewPaymentRepository(dao dao.PaymentDAO) PaymentRepository {
	return &paymentRepository{
		dao: dao,
	}
}

func (repo *paymentRepository) AddPayment(ctx context.Context, pmt domain.Payment) error {
	return repo.dao.Insert(ctx, repo.toEntity(pmt))
}

func (repo *paymentRepository) UpdatePayment(ctx context.Context, pmt domain.Payment) error {
	return repo.dao.UpdateTxnIDAndStatus(ctx, pmt.BizTradeNO, pmt.TxnID, pmt.Status)
}

func (repo *paymentRepository) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]domain.Payment, error) {
	res, err := repo.dao.FindExpiredPayment(ctx, offset, limit, t)
	payments := make([]domain.Payment, len(res))
	for i, p := range res {
		payments[i] = repo.toDomain(p)
	}
	return payments, err
}

func (repo *paymentRepository) GetPayment(ctx context.Context, bizTradeNO string) (domain.Payment, error) {
	payment, err := repo.dao.GetPayment(ctx, bizTradeNO)
	return repo.toDomain(payment), err
}

func (repo *paymentRepository) toDomain(payment dao.Payment) domain.Payment {
	return domain.Payment{
		Amt: domain.Amount{
			Total:    payment.Amt,
			Currency: payment.Currency,
		},
		BizTradeNO:  payment.BizTradeNO,
		Description: payment.Description,
		Status:      domain.PaymentStatus(payment.Status),
		TxnID:       payment.TxnID.String,
	}
}

func (repo *paymentRepository) toEntity(payment domain.Payment) dao.Payment {
	return dao.Payment{
		Amt:         payment.Amt.Total,
		Currency:    payment.Amt.Currency,
		BizTradeNO:  payment.BizTradeNO,
		Description: payment.Description,
		Status:      payment.Status.AsUint8(),
		TxnID: sql.NullString{
			String: payment.TxnID,
			Valid:  payment.TxnID != "",
		},
	}
}
