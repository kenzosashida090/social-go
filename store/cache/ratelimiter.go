package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiterRedisStorage struct {
	rdb *redis.Client
}

func (c *RateLimiterRedisStorage) Get(ctx context.Context, addrs string) (int64, error) {

	data, err := c.rdb.Get(ctx, addrs).Result()
	if err == redis.Nil {
		fmt.Println("redis empty", err)
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	fmt.Println(data, "CACHE RATE LIMITER")
	var count int64
	err = json.Unmarshal([]byte(data), &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
func (c *RateLimiterRedisStorage) Set(ctx context.Context, addrs string, count int64) error {
	err := c.rdb.Set(ctx, addrs, count, time.Second*60).Err()
	fmt.Println("----------------REDIS")
	if err != nil {
		return err
	}

	return nil

}
