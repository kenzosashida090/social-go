package cache

import (
	"context"

	"github.com/kenzosashida090/social/store"
	"github.com/redis/go-redis/v9"
)

type RateLimiterType struct {
	Addrs string `json:"addrs"`
	Count int64  `json:"count"`
}
type UserStorage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}
type RateLimitStorage struct {
	RateLimiter interface {
		Get(context.Context, string) (int64, error)
		Set(context.Context, string, int64) error
	}
}

func NewRedisStorage(redisClient *redis.Client) UserStorage {
	return UserStorage{
		Users: &UserRedisStorage{rdb: redisClient},
	}
}

func NewRedisLimitRateStorage(redisClient *redis.Client) RateLimitStorage {
	return RateLimitStorage{
		RateLimiter: &RateLimiterRedisStorage{rdb: redisClient},
	}
}
