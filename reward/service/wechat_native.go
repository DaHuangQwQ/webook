package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	accountv1 "webook/api/proto/gen/account/v1"
	paymentv1 "webook/api/proto/gen/payment/v1"
	"webook/pkg/logger"
	"webook/reward/domain"
	"webook/reward/repository"
)

type WechatNativeRewardService struct {
	client        paymentv1.WechatPaymentServiceClient
	accountClient accountv1.AccountServiceClient
	repo          repository.RewardRepository
	l             logger.LoggerV1
}

func NewWechatNativeRewardService(client paymentv1.WechatPaymentServiceClient, accountClient accountv1.AccountServiceClient, repo repository.RewardRepository, l logger.LoggerV1) RewardService {
	return &WechatNativeRewardService{
		client:        client,
		accountClient: accountClient,
		repo:          repo,
		l:             l,
	}
}

func (svc *WechatNativeRewardService) PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	// 缓存
	res, err := svc.repo.GetCachedCodeURL(ctx, r)
	if err == nil {
		return res, nil
	}
	r.Status = domain.RewardStatusInit
	rid, err := svc.repo.CreateReward(ctx, r)
	if err != nil {
		return domain.CodeURL{}, err
	}
	paymentRes, err := svc.client.NativePrePay(ctx, &paymentv1.PrePayRequest{
		Amt: &paymentv1.Amount{
			Total:    r.Amt,
			Currency: "CNY",
		},
		BizTradeNo:  svc.bizTradeNO(rid),
		Description: fmt.Sprintf("打赏-%s", r.Target.BizName),
	})
	if err != nil {
		return domain.CodeURL{}, err
	}
	codeUrl := domain.CodeURL{
		Rid: r.Id,
		URL: paymentRes.CodeUrl,
	}
	// 缓存
	err = svc.repo.CachedCodeURL(ctx, codeUrl, r)
	return codeUrl, err
}

func (svc *WechatNativeRewardService) GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error) {
	// 快路径
	res, err := svc.repo.GetReward(ctx, rid)
	if err != nil {
		return domain.Reward{}, err
	}
	// 确保是自己打赏
	if res.Uid != uid {
		return domain.Reward{}, errors.New("非法访问别人的打赏记录")
	}
	// 降级或者限流的时候，不走慢路径
	if ctx.Value("limited") == "true" {
		return res, nil
	}
	if !res.Completed() {
		// 我去问一下，有可能支付那边已经处理好了，已经收到回调了
		pmtRes, err := svc.client.GetPayment(ctx, &paymentv1.GetPaymentRequest{
			BizTradeNo: svc.bizTradeNO(rid),
		})
		if err != nil {
			svc.l.Error("慢路径查询支付状态失败",
				logger.Error(err),
				logger.Int64("rid", rid))
			return res, nil
		}
		switch pmtRes.Status {
		case paymentv1.PaymentStatus_PaymentStatusSuccess:
			res.Status = domain.RewardStatusPayed
		case paymentv1.PaymentStatus_PaymentStatusInit:
			res.Status = domain.RewardStatusInit
		case paymentv1.PaymentStatus_PaymentStatusRefund:
			res.Status = domain.RewardStatusFailed
		case paymentv1.PaymentStatus_PaymentStatusFailed:
			res.Status = domain.RewardStatusFailed
		case paymentv1.PaymentStatus_PaymentStatusUnknown:
		}
		err = svc.UpdateReward(ctx, svc.bizTradeNO(rid), res.Status)
		if err != nil {
			svc.l.Error("慢路径更新本地状态失败", logger.Error(err),
				logger.Int64("rid", rid))
		}
	}
	return res, nil
}

func (svc *WechatNativeRewardService) UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error {
	rid := svc.toRid(bizTradeNO)
	err := svc.repo.UpdateStatus(ctx, rid, status)
	if err != nil {
		return err
	}
	// 完成了支付，准备入账
	if status == domain.RewardStatusPayed {
		r, err := svc.repo.GetReward(ctx, rid)
		if err != nil {
			return err
		}
		// webook 抽成
		weAmt := int64(float64(r.Amt) * 0.1)
		_, err = svc.accountClient.Credit(ctx, &accountv1.CreditRequest{
			Biz:   "reward",
			BizId: rid,
			Items: []*accountv1.CreditItem{
				{
					AccountType: accountv1.AccountType_AccountTypeReward,
					// 虽然可能为 0，但是也要记录出来
					Amt:      weAmt,
					Currency: "CNY",
				},
				{
					Account:     r.Uid,
					Uid:         r.Uid,
					AccountType: accountv1.AccountType_AccountTypeReward,
					Amt:         r.Amt - weAmt,
					Currency:    "CNY",
				},
			},
		})
		if err != nil {
			svc.l.Error("入账失败了，快来修数据啊！！！",
				logger.String("biz_trade_no", bizTradeNO),
				logger.Error(err))
			// 做好监控和告警，这里
			return err
		}
	}
	return nil
}

func (svc *WechatNativeRewardService) bizTradeNO(rid int64) string {
	return fmt.Sprintf("reward-%d", rid)
}

func (svc *WechatNativeRewardService) toRid(tradeNO string) int64 {
	ridStr := strings.Split(tradeNO, "-")
	val, _ := strconv.ParseInt(ridStr[1], 10, 64)
	return val
}
