package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"webook/sms/service"
)

type SMSService struct {
	svc service.Service
	key []byte
}

func (s *SMSService) Send(ctx context.Context, tplToken string, args []service.NamedArg, numbers ...string) error {
	var claims SMSClaims
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.svc.Send(ctx, claims.Tpl, args, numbers...)
}

type SMSClaims struct {
	jwt.RegisteredClaims
	Tpl string
	// 额外加字段
}
