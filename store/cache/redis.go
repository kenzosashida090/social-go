package cache

import (
	"github.com/redis/go-redis/v9"
)

type RedisClientConn struct {
}

func NewConnectionClient(addrs string, pass string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addrs,
		Password: pass,
		DB:       db,
		Protocol: 2,
	})

	return rdb

}
