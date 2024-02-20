package auth

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc sms.Service
	key string
}

func (S SMSService) Send(ctx context.Context, tplToken string, params []string, phoneNumbers ...string) error {
	var claims SMSClaims
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return S.key, nil
	})
	if err != nil {
		return err
	}
	return S.svc.Send(ctx, claims.Tpl, params, phoneNumbers...)
}

type SMSClaims struct {
	jwt.RegisteredClaims
	Tpl string
}
