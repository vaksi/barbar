package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	parentKey string
	exp       time.Duration
	conn      *redis.Client
}

type Config struct {
	Connection string
	Password   string
	DB         int
	Expiration time.Duration
}

var ErrRedisKeyNil = errors.New("data_nil")

func NewRedis(parentKey string, redisCfg Config) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Connection,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &Redis{
		parentKey: parentKey,
		conn:      rdb,
	}, nil
}

func (r *Redis) SetDataRedis(ctx context.Context, key string, body interface{}) error {
	key = r.parentKey + "-" + key
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = r.conn.Set(ctx, key, buf, r.exp).Err()
	return err
}

func (r *Redis) GetDataRedis(ctx context.Context, key string, body interface{}) error {
	key = r.parentKey + "-" + key
	buf, err := r.conn.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrRedisKeyNil
	} else if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(buf), &body)

	return err
}

func (r *Redis) DeleteDataRedis(ctx context.Context, key string) error {
	key = r.parentKey + "-" + key

	err := r.conn.Del(ctx, key).Err()
	if err == redis.Nil {
		return ErrRedisKeyNil
	} else if err != nil {
		return err
	}
	return nil
}
