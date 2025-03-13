package cache

import (
	"github.com/davycun/eta/pkg/common/caller"
	"strings"
	"time"
)

func Set(key string, value interface{}) error {
	return defaultCacheRedis.Set(key, value)
}

func SetEx(key string, value interface{}, expiration time.Duration) error {
	return defaultCacheRedis.SetEx(key, value, expiration)
}
func Exists(key string) (error, bool) {
	return defaultCacheRedis.Exists(key)
}
func TTL(key string) (error, time.Duration) {
	return defaultCacheRedis.TTL(key)
}
func Expire(key string, expiration time.Duration) (error, bool) {
	return defaultCacheRedis.Expire(key, expiration)
}
func Get(key string, dest any) (error, bool) {
	return defaultCacheRedis.Get(key, dest)
}
func Detail(key string) (error, *any, *time.Duration) {
	return defaultCacheRedis.Detail(key)
}

func Del(key ...string) (error, bool) {
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

func DelKeyPattern(key string) (err error, b bool) {
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
				err, _ = defaultCacheRedis.Del(k)
				if err != nil {
					return err
				}
			}
			return nil
		}).Err
	return err, err == nil
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
