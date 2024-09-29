package grpc

import (
	"context"
	"github.com/DaHuangQwQ/gutil/slice"
	"google.golang.org/grpc"
	"webook/account/domain"
	"webook/account/service"
	accountv1 "webook/api/proto/gen/account/v1"
)

type AccountServiceServer struct {
	accountv1.UnimplementedAccountServiceServer
	svc service.AccountService
}

func NewAccountServiceServer(svc service.AccountService) *AccountServiceServer {
	return &AccountServiceServer{
		svc: svc,
	}
}

func (s *AccountServiceServer) Register(server *grpc.Server) {
	accountv1.RegisterAccountServiceServer(server, s)
}

func (s *AccountServiceServer) Credit(ctx context.Context, request *accountv1.CreditRequest) (*accountv1.CreditResponse, error) {
	err := s.svc.Credit(ctx, s.toDomain(request))
	return &accountv1.CreditResponse{}, err
}

func (s *AccountServiceServer) toDomain(c *accountv1.CreditRequest) domain.Credit {
	return domain.Credit{
		Biz:   c.Biz,
		BizId: c.BizId,
		Items: slice.Map(c.Items, func(idx int, src *accountv1.CreditItem) domain.CreditItem {
			return s.itemToDomain(src)
		}),
	}
}

func (s *AccountServiceServer) itemToDomain(c *accountv1.CreditItem) domain.CreditItem {
	return domain.CreditItem{
		Account: c.Account,
		Amt:     c.Amt,
		Uid:     c.Uid,
		// 两者取值都是一样的，我偷个懒，直接转
		AccountType: domain.AccountType(c.AccountType),
		Currency:    c.Currency,
	}
}
