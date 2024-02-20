package tencent

import (
	"basic-go/webook/pkg/limiter"
	"context"
	"errors"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	"go.uber.org/zap"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client    *sms.Client
	appId     *string
	signature *string
	limiter   limiter.Limiter
}

func (s *Service) Send(ctx context.Context, tplId string, params []string, phoneNumbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "tecent-sms-service")
	if err != nil {
		return err
	}
	if limited {
		return errors.New("触发了限流")
	}
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = s.appId
	request.SignName = s.signature
	request.TemplateId = ekit.ToPtr(tplId)
	request.TemplateParamSet = s.toPtrSlice(params)
	request.PhoneNumberSet = s.toPtrSlice(phoneNumbers)

	response, err := s.client.SendSms(request)
	zap.L().Debug("请求腾讯  SendSMS 接口", zap.Any("rep", request), zap.Any("resp", response))
	// 处理异常
	if err != nil {
		return err
	}

	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			continue
		}
		status := *statusPtr
		if status.Code == nil || *status.Code != "Ok" {
			return fmt.Errorf("发送短信失败 code: %s, msg: %s", *status.Code, *status.Message)
		}
	}
	return nil

}

func (s *Service) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data, func(idx int, src string) *string {
		return &src
	})
}

func NewService(client *sms.Client, appId string, signature string, l limiter.Limiter) *Service {
	return &Service{
		client:    client,
		appId:     ekit.ToPtr(appId),
		signature: ekit.ToPtr(signature),
		limiter:   l,
	}
}
