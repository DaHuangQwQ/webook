package grpc

import (
	"context"
	"google.golang.org/grpc"
	paymentv1 "webook/api/proto/gen/payment/v1"
	"webook/payment/domain"
	"webook/payment/service/wechat"
)

type PaymentServiceServer struct {
	paymentv1.UnimplementedWechatPaymentServiceServer
	svc *wechat.NativePaymentService
}

func NewPaymentServiceServer(svc *wechat.NativePaymentService) *PaymentServiceServer {
	return &PaymentServiceServer{
		svc: svc,
	}
}

func (p *PaymentServiceServer) Register(server *grpc.Server) {
	paymentv1.RegisterWechatPaymentServiceServer(server, p)
}

func (p *PaymentServiceServer) NativePrePay(ctx context.Context, request *paymentv1.PrePayRequest) (*paymentv1.NativePrePayResponse, error) {
	codeUrl, err := p.svc.Prepay(ctx, domain.Payment{
		Amt: domain.Amount{
			Total:    request.Amt.Total,
			Currency: request.Amt.Currency,
		},
		BizTradeNO:  request.BizTradeNo,
		Description: request.Description,
	})
	return &paymentv1.NativePrePayResponse{
		CodeUrl: codeUrl,
	}, err
}

func (p *PaymentServiceServer) GetPayment(ctx context.Context, request *paymentv1.GetPaymentRequest) (*paymentv1.GetPaymentResponse, error) {
	payment, err := p.svc.GetPayment(ctx, request.GetBizTradeNo())
	return &paymentv1.GetPaymentResponse{
		Status: paymentv1.PaymentStatus(payment.Status),
	}, err
}
