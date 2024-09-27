package wechat

import (
	"context"
	"errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"time"
	"webook/payment/domain"
	"webook/payment/events"
	"webook/payment/repository"
	"webook/pkg/logger"
)

type NativePaymentService struct {
	appID                string
	mchID                string
	notifyURL            string
	repo                 repository.PaymentRepository
	svc                  *native.NativeApiService
	producer             events.Producer
	l                    logger.LoggerV1
	nativeCBTypeToStatus map[string]domain.PaymentStatus
}

func NewNativePaymentService(
	appID string, mchID string,
	repo repository.PaymentRepository,
	svc *native.NativeApiService,
	producer events.Producer,
	l logger.LoggerV1,
) *NativePaymentService {
	return &NativePaymentService{
		appID: appID,
		mchID: mchID,
		// 从配置文件读取 测试环境和线上环境
		notifyURL: "http://wechat.dahuang.pro/pay/callback",
		repo:      repo,
		svc:       svc,
		producer:  producer,
		l:         l,
		nativeCBTypeToStatus: map[string]domain.PaymentStatus{
			"SUCCESS":    domain.PaymentStatusSuccess,
			"PAYERROR":   domain.PaymentStatusFailed,
			"NOTPAY":     domain.PaymentStatusInit,
			"USERPAYING": domain.PaymentStatusInit,
			"CLOSED":     domain.PaymentStatusFailed,
			"REVOKED":    domain.PaymentStatusFailed,
			"REFUND":     domain.PaymentStatusRefund,
			// 其它状态都可以加
		},
	}
}

func (svc *NativePaymentService) Prepay(ctx context.Context, pmt domain.Payment) (string, error) {
	pmt.Status = domain.PaymentStatusInit
	err := svc.repo.AddPayment(ctx, pmt)
	if err != nil {
		return "", err
	}
	//sn := uuid.New().String()
	resp, _, err := svc.svc.Prepay(ctx, native.PrepayRequest{
		Appid:       core.String(svc.appID),
		Mchid:       core.String(svc.mchID),
		Description: core.String(pmt.Description),
		OutTradeNo:  core.String(pmt.BizTradeNO),
		// 最好这个要带上
		TimeExpire: core.Time(time.Now().Add(time.Minute * 30)),
		Amount: &native.Amount{
			Total:    core.Int64(pmt.Amt.Total),
			Currency: core.String(pmt.Amt.Currency),
		},
	})

	if err != nil {
		return "", err
	}
	return *resp.CodeUrl, nil
}

func (svc *NativePaymentService) HandlerCallback(ctx context.Context, txn *payments.Transaction) error {
	return svc.updateByTxn(ctx, txn)
}

// SyncWechatInfo 兜底，同步状态 启动定时任务
func (svc *NativePaymentService) SyncWechatInfo(ctx context.Context, bizTradeNo string) error {
	txn, _, err := svc.svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(bizTradeNo),
		Mchid:      core.String(svc.mchID),
	})
	if err != nil {
		return err
	}
	return svc.updateByTxn(ctx, txn)
}

func (svc *NativePaymentService) updateByTxn(ctx context.Context, txn *payments.Transaction) error {
	status, ok := svc.nativeCBTypeToStatus[*txn.TradeState]
	if !ok {
		svc.l.Warn("微信支付状态映射失败")
		return errors.New("微信支付状态映射失败")
	}
	return svc.repo.UpdatePayment(ctx, domain.Payment{
		BizTradeNO: *txn.OutTradeNo,
		Status:     status,
		TxnID:      *txn.TransactionId,
	})
}
