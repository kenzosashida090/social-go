package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/kenzosashida090/social/store"
	"github.com/redis/go-redis/v9"
)

type UserRedisStorage struct {
	rdb *redis.Client
}
type UserCache struct {
	User *store.User `redis:"user"`
}

const UserTimeExp = time.Minute

func (c *UserRedisStorage) Get(ctx context.Context, userId int64) (*store.User, error) {
	var user store.User

	val := strconv.FormatInt(userId, 10)
	data, err := c.rdb.Get(ctx, val).Result()
	if err == redis.Nil {
		fmt.Println("redis empty", err)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	fmt.Println(user, "CACHE User")
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (c *UserRedisStorage) Set(ctx context.Context, user *store.User) error {

	userId := strconv.FormatInt(user.ID, 10)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = c.rdb.Set(ctx, userId, data, UserTimeExp).Err()
	fmt.Println(data, "----------------REDIS")
	if err != nil {
		return err
	}

	return nil

}
