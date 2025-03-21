package cache

import (
	"github.com/davycun/eta/pkg/common/caller"
	"strings"
	"time"
)

func Set(key string, value interface{}) error {
	return defaultCacheRedis.Set(key, value)
}

// SetEx
// expiration 单位为纳秒，存入redis的时候，会转换为秒即expiration/time.Second
// 本地存储不转换只是在当前时间的基础上加上expiration，其实单位也是纳秒
func SetEx(key string, value interface{}, expiration time.Duration) error {
	return defaultCacheRedis.SetEx(key, value, expiration)
}
func Exists(key string) (bool, error) {
	return defaultCacheRedis.Exists(key)
}
func TTL(key string) (time.Duration, error) {
	return defaultCacheRedis.TTL(key)
}
func Expire(key string, expiration time.Duration) (error, bool) {
	return defaultCacheRedis.Expire(key, expiration)
}
func Get(key string, dest any) (bool, error) {
	return defaultCacheRedis.Get(key, dest)
}
func Detail(key string) (any, time.Duration, error) {
	return defaultCacheRedis.Detail(key)
}

func Del(key ...string) (bool, error) {
	return defaultCacheRedis.Del(key...)
}

func AddAfterDel(f ...AfterDel) {
	for i, _ := range f {
		defaultCacheRedis.AddAfterDel(f[i])
	}
}
func PublishDelKey(key ...string) error {
	return defaultCacheRedis.PublishDelKey(key...)
}

func DelKeyPattern(key string) (b bool, err error) {
	if !strings.HasSuffix(key, "*") {
		return Del(key)
	}
	keys := make([]string, 0)
	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			keys, err = defaultCacheRedis.Keys(key)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			for _, k := range keys {
				_, err = defaultCacheRedis.Del(k)
				if err != nil {
					return err
				}
			}
			return nil
		}).Err
	return err == nil, err
}

func Scan(cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error) {
	if match == "" {
		match = "*"
	}
	if count > 10000 || count <= 0 {
		count = 10000
	}
	return defaultCacheRedis.Scan(cursor, match, count)
}

func Unlink(keys ...string) error {
	return defaultCacheRedis.Unlink(keys...)
}
