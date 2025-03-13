package cache

import (
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	defaultCacheRedis *redisCache
	defaultEmptyCache *cacheAdapter
	once              sync.Once
)

func init() {
	defaultEmptyCache = &cacheAdapter{}
	defaultCacheRedis = &redisCache{}
}

// InitCacheWithRedis
// 减少对global的依赖，暂时先通过server自动复制
func InitCacheWithRedis(rds *redis.Client) {
	//貌似锁不锁都无所谓
	once.Do(func() {
		defaultCacheRedis.rds = rds
		go defaultCacheRedis.subscribeDeleteCommand()
	})
}

func GetCache(tp Type, option *Option) Cache {
	switch tp {
	case Redis:
		return defaultCacheRedis
	default:
		return defaultEmptyCache
	}
}
