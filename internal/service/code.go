package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

//type CodeService interface {
//}

type CodeService struct {
	repo   *repository.CodeRepository
	smsMvc sms.Service
}

const (
	codeTplId = "123456"
)

func NewCodeService(repo *repository.CodeRepository, smsMvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsMvc: smsMvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return svc.smsMvc.Send(ctx, codeTplId, []sms.NamedArg{{Value: code}}, phone)
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
