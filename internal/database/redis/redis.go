package redis

import (
	"github.com/go-redis/redis"
	"test_repository/internal/config"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})
}
