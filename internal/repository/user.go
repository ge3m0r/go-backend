package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
	"database/sql"
)

var (
	ErrDuplicateUser = dao.ErrDuplicated
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	EditProfile(c context.Context, user domain.User) error
	Profile(c context.Context, user domain.User) (domain.User, error)
	FindByPhone(c context.Context, phone string) (domain.User, error)
	FindByWechat(ctx context.Context, openId string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindbyEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), err
}

func (repo *CacheUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		ID:       u.ID,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		NickName: u.NickName,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday,
		WechatInfo: domain.WechatInfo{
			OpenId:  u.WechatOpenId.String,
			UnionId: u.WechatUnionId.String,
		},
	}

}

func (repo *CacheUserRepository) EditProfile(c context.Context, user domain.User) error {
	err := repo.dao.EditProfile(c, user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *CacheUserRepository) Profile(c context.Context, user domain.User) (domain.User, error) {
	uc, err := repo.cache.Get(c, user.ID)
	if err == nil {
		return uc, err
	}
	u, err := repo.dao.Profile(c, user)
	if err != nil {
		return domain.User{}, err
	}
	du := repo.toDomain(u)
	err = repo.cache.Set(c, du)
	return du, nil
}

func (repo *CacheUserRepository) FindByPhone(c context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(c, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), err
}

func (repo *CacheUserRepository) FindByWechat(ctx context.Context, openId string) (domain.User, error) {
	ue, err := repo.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(ue), nil
}

func (repo *CacheUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		ID: u.ID,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		NickName: u.NickName,
		Ctime:    u.Ctime,
		Utime:    u.Utime,
		WechatUnionId: sql.NullString{
			String: u.WechatInfo.UnionId,
			Valid:  u.WechatInfo.UnionId != "",
		},
		WechatOpenId: sql.NullString{
			String: u.WechatInfo.OpenId,
			Valid:  u.WechatInfo.OpenId != "",
		},
	}

}
