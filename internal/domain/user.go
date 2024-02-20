package domain

type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement:true"`
	Email    string `gorm:"unique"`
	Password string
	NickName string
	AboutMe  string
	Birthday string
	Phone    string
	Ctime    int64
	Utime    int64

	WechatInfo WechatInfo `gorm:"embedded"`
}

func (u User) ValidateEmail() string {
	return u.Email
}
