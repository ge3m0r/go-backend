package ioc

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/localsms"
)

func InitSMSService() sms.Service {
	return localsms.NewService()
}

func initTencentSmsService() sms.Service {
	return nil
}
