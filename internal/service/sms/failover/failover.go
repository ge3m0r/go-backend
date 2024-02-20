package failover

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"log"
)

type FailOverSMSService struct {
	svcs []sms.Service

	idx uint64
}

func (f FailOverSMSService) Send(ctx context.Context, tplId string, params []string, phoneNumbers ...string) error {
	//TODO implement me
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, params, phoneNumbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("轮询了所有的服务器但是都失败了")
}

func (f FailOverSMSService) SendV1(ctx context.Context, tplId string, params []string, phoneNumbers ...string) error {
	idx := f.idx
	length := uint64(len(f.svcs))

	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, params, phoneNumbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		}

	}
	return errors.New("轮询了所有的服务器但是都失败了")
}
