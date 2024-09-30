package service

import (
	"context"
	"fmt"
	"math/rand"
	smsv1 "webook/api/proto/gen/sms/v1"
	"webook/code/repository"
)

type codeService struct {
	repo   repository.CodeRepository
	smsMvc smsv1.SmsServiceClient
}

const (
	codeTplId = "SMS_472455063"
)

func NewCodeService(repo repository.CodeRepository, smsMvc smsv1.SmsServiceClient) CodeService {
	return &codeService{
		repo:   repo,
		smsMvc: smsMvc,
	}
}

func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	_, err = svc.smsMvc.Send(ctx, &smsv1.SmsSendRequest{
		TplId: codeTplId,
		Args: []*smsv1.NameArgs{
			{
				Name:  "code",
				Value: code,
			},
		},
		Numbers: []string{phone},
	})
	return err
}

func (svc *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *codeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
