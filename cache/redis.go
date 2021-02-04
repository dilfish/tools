package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	Ctx context.Context
	Cli *redis.Client
	Opt *redis.Options
}

func NewRedis(opt *redis.Options) (*RedisService, error) {
	rdb := redis.NewClient(opt)
	return &RedisService{Ctx: context.Background(), Cli: rdb, Opt: opt}, nil
}

func (r *RedisService) Get(key string) (string, error) {
	val, err := r.Cli.Get(r.Ctx, key).Result()
	return val, err
}

func (r *RedisService) Set(key, val string) error {
	return r.Cli.Set(r.Ctx, key, val, 0).Err()
}

func (r *RedisService) IsNil(e interface{}) bool {
	err, ok := e.(error)
	if !ok {
		return false
	}
	if err == redis.Nil {
		return true
	}
	return false
}
