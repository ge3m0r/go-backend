package sms

import "context"

type Service interface {
	Send(ctx context.Context, tplId string, params []string, phoneNumbers ...string) error
}
