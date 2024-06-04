package cache

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

type RedisInstance struct {
	Pool *redis.Pool
}

func NewRedisInstance(pool *redis.Pool) *RedisInstance {
	return &RedisInstance{Pool: pool}
}

func (r *RedisInstance) Get(ctx context.Context, key string) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	data, err := redis.Int(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}

	return data, nil
}

func (r *RedisInstance) Incr(ctx context.Context, key string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("INCR", key)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisInstance) Expire(ctx context.Context, key string, ttlSecs int) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, ttlSecs)
	if err != nil {
		return err
	}
	return nil
}
