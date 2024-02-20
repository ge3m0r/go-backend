package dao

import (
	"basic-go/webook/internal/domain"
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicated     = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindbyEmail(ctx context.Context, email string) (User, error)
	EditProfile(c context.Context, user domain.User) error
	Profile(c context.Context, user domain.User) (User, error)
	FindByPhone(c context.Context, phone string) (User, error)
	FindByWechat(c context.Context, phone string) (User, error)
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GormUserDAO{
		db: db,
	}
}

func (dao *GormUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	println(err)
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicate uint16 = 1062
		if me.Number == duplicate {
			return ErrDuplicated
		}
	}
	return err
}

func (dao *GormUserDAO) FindByWechat(ctx context.Context, openId string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id=?", openId).First(&u).Error
	return u, err
}

func (dao *GormUserDAO) FindbyEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *GormUserDAO) EditProfile(c context.Context, user domain.User) error {
	var u = &User{
		NickName: user.NickName,
		Birthday: user.Birthday,
		AboutMe:  user.AboutMe,
	}

	err := dao.db.WithContext(c).Model(&u).Where("id=?", u.ID).Updates(map[string]any{
		"nick_name": u.NickName,
		"birthday":  u.Birthday,
		"about_me":  u.AboutMe,
	}).Error
	return err
}

func (dao *GormUserDAO) Profile(c context.Context, user domain.User) (User, error) {
	var u = &User{

		ID: user.ID,
	}
	err := dao.db.WithContext(c).Select("nick_name", "birthday", "about_me").Where("id=?", u.ID).Find(&u).Error
	return *u, err
}

func (dao *GormUserDAO) FindByPhone(c context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(c).Where("phone=?", phone).First(&u).Error
	return u, err
}

type User struct {
	ID       int64          `gorm:"primaryKey, autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	NickName string `gorm:"type=varchar(128)"`
	Birthday string
	Phone    sql.NullString `gorm:"unique"`
	AboutMe  string         `gorm:"type=varchar(4096)"`

	Ctime         int64
	Utime         int64
	WechatOpenId  sql.NullString `gorm:"unique"`
	WechatUnionId sql.NullString
}
