package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	cachemocks "basic-go/webook/internal/repository/cache/mocks"
	"basic-go/webook/internal/repository/dao"
	daomocks "basic-go/webook/internal/repository/dao/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheUserRepository_EditProfile(t *testing.T) {
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO)
		ctx       context.Context
		givenUser domain.User
		wantErr   error
	}{
		{
			name: "查找成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				user := domain.User{
					ID: 123,
				}
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), user.ID).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().EditProfile(gomock.Any(), user).Return(nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					ID:       123,
					Email:    "123@qq.com",
					Password: "123456",
					Birthday: "2012-12-23",
					AboutMe:  "自我介绍",
					Phone:    "13125170185",
					Ctime:    time.Now().UnixMilli(),
				}).Return(nil)
				return c, d
			},
			ctx: context.Background(),
			givenUser: domain.User{
				ID: 123,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc, ud := tc.mock(ctrl)
			svc := NewUserRepository(ud, uc)
			err := svc.EditProfile(tc.ctx, tc.givenUser)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
