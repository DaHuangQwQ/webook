package localsms

import (
	"context"
	"log"
	"webook/internal/service/sms"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tpl string, args []sms.NamedArg, numbers ...string) error {
	log.Println("验证码是", args)
	return nil
}
