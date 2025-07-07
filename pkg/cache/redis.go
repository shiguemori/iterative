package redis

import (
	"context"
	"fmt"
	"time"

	"b3-challenge/internal/config"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisCache struct {
	redis *redis.Client
}

func NewRedisClient(cfg *config.RedisConfig) *RedisCache {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &RedisCache{
		redis: rdb,
	}
}

func (r *RedisCache) SetJSON(key string, value interface{}, ttl time.Duration) error {
	return r.redis.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) GetString(key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}
