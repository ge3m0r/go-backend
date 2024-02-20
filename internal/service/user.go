package service

import (
	"basic-go/webook/internal/domain"
	repository "basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"context"
	"errors"
	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUser         = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("账号或密码不存在")
	ErrCodeSendTooMany       = cache.ErrSendTooMany
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	Edit(c context.Context, user domain.User) error
	Profile(c context.Context, user domain.User) (domain.User, error)
	FindOrCreate(c context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, err
}

func (svc *userService) Edit(c context.Context, user domain.User) error {
	err := svc.repo.EditProfile(c, user)
	if err != nil {
		return err
	}
	return nil

}

func (svc *userService) Profile(c context.Context, user domain.User) (domain.User, error) {
	u, err := svc.repo.Profile(c, user)
	return u, err
}

func (svc *userService) FindOrCreate(c context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(c, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	err = svc.repo.Create(c, domain.User{
		Phone: phone,
	})
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	return svc.repo.FindByPhone(c, phone)
}

func (svc *userService) FindOrCreateByWechat(c context.Context, wechatInfo domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(c, wechatInfo.OpenId)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	zap.L().Info("新用户", zap.Any("wechatinfo", wechatInfo))
	err = svc.repo.Create(c, domain.User{
		WechatInfo: wechatInfo,
	})
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	return svc.repo.FindByWechat(c, wechatInfo.OpenId)
}
