package localsms

import (
	"context"
	"log"
	"webook/sms/service"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tpl string, args []service.NamedArg, numbers ...string) error {
	log.Println("验证码是", args)
	return nil
}
