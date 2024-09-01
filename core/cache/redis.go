package cache

import (
	"context"

	"github.com/nineteenseventy/minichat/core/logging"
	"github.com/redis/go-redis/v9"
)

var globalRedis *redis.Client

func InitRedis(ctx context.Context, config redis.Options) error {
	logger := logging.GetLogger("redis")

	globalRedis = redis.NewClient(&config)
	info, err := globalRedis.InfoMap(ctx, "server").Result()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to Redis")
		return err
	}

	host := globalRedis.Options().Addr
	redisVersion := info["Server"]["redis_version"]
	logger.Info().Str("version", redisVersion).Str("host", host).Msg("Connected to Redis")
	return nil
}

func GetRedis() *redis.Client {
	if globalRedis == nil {
		panic("Redis not initialized")
	}
	return globalRedis
}
