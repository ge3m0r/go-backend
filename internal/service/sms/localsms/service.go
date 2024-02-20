package localsms

import "basic-go/webook/internal/service/sms"

type Service struct {
	sms.Service
}

func NewService() *Service {
	return &Service{}
}
