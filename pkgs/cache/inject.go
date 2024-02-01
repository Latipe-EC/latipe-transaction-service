package cache

import (
	"fmt"
	"github.com/google/wire"
	"latipe-transaction-service/config"
	"latipe-transaction-service/pkgs/cache/redis"
	"time"
)

var Set = wire.NewSet(
	NewCacheEngine,
)

func NewCacheEngine(config *config.Config) (*redisCache.CacheEngine, error) {
	cfg := redisCache.RedisConfig{
		Address:               fmt.Sprintf("%v:%v", config.Cache.Redis.Address, config.Cache.Redis.Port),
		DB:                    config.Cache.Redis.DB,
		Password:              config.Cache.Redis.Password,
		ContextTimeoutEnabled: true,
		PoolSize:              5,
		PoolTimeout:           5,
		DialTimeout:           5,
		ReadTimeout:           5 * time.Second,
		WriteTimeout:          5 * time.Second,
		ConnectTimeout:        5 * time.Second,
	}
	client, err := redisCache.NewCacheEngine(cfg)
	if err != nil {
		return nil, err
	}
	return client, err
}
