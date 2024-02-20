package cache

import (
	"basic-go/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(c context.Context, id int64) (domain.User, error)
	Set(c context.Context, du domain.User) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (uc *RedisUserCache) Get(c context.Context, id int64) (domain.User, error) {
	key := uc.key(id)
	data, err := uc.cmd.Get(c, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (c *RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid, uid)
}

func (uc *RedisUserCache) Set(c context.Context, du domain.User) error {
	key := uc.key(du.ID)
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return uc.cmd.Set(c, key, data, uc.expiration).Err()
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
