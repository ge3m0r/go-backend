package ratelimit

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/limiter"
	"context"
	"errors"
)

var errLimited = errors.New("被限流了")

type RateLimitSMSService struct {
	svc   sms.Service
	limit limiter.Limiter
	key   string
}

func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, phoneNumbers ...string) error {
	Limited, err := r.limit.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if Limited {
		return errLimited
	}
	return r.Send(ctx, tplId, args, phoneNumbers...)
}

func NewRateLimitSMSService(svc sms.Service, limiter limiter.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:   svc,
		limit: limiter,
		key:   "sms-limiter",
	}

}
