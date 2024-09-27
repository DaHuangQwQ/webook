package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
	"webook/payment/domain"
)

type PaymentGORMDAO struct {
	db *gorm.DB
}

func NewPaymentGORMDAO(db *gorm.DB) PaymentDAO {
	return &PaymentGORMDAO{db: db}
}

func (dao *PaymentGORMDAO) Insert(ctx context.Context, pmt Payment) error {
	//TODO implement me
	panic("implement me")
}

func (dao *PaymentGORMDAO) UpdateTxnIDAndStatus(ctx context.Context, bizTradeNo string, txnID string, status domain.PaymentStatus) error {
	//TODO implement me
	panic("implement me")
}

func (dao *PaymentGORMDAO) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]Payment, error) {
	var res []Payment
	err := dao.db.WithContext(ctx).Where("status = ? AND u_time < ?", uint8(domain.PaymentStatusInit), t.UnixMilli()).
		Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (dao *PaymentGORMDAO) GetPayment(ctx context.Context, bizTradeNO string) (Payment, error) {
	//TODO implement me
	panic("implement me")
}
