package repository

import (
	"context"
	"time"
	"webook/account/domain"
	"webook/account/repository/dao"
)

type accountRepository struct {
	dao dao.AccountDAO
}

func NewAccountRepository(dao dao.AccountDAO) AccountRepository {
	return &accountRepository{
		dao: dao,
	}
}

func (a *accountRepository) AddCredit(ctx context.Context, c domain.Credit) error {
	activities := make([]dao.AccountActivity, 0)
	now := time.Now().UnixMilli()
	for _, itm := range c.Items {
		activities = append(activities, dao.AccountActivity{
			Uid:         itm.Uid,
			Biz:         c.Biz,
			BizId:       c.BizId,
			Account:     itm.Account,
			AccountType: itm.AccountType.AsUint8(),
			Amount:      itm.Amt,
			Currency:    itm.Currency,
			CTime:       now,
			UTime:       now,
		})
	}
	return a.dao.AddActivities(ctx, activities...)
}
