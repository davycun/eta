package cache

import (
	"github.com/davycun/eta/pkg/common/logger"
	"time"
)

const (
	Redis Type = "redis"
)

type Type string

type Option struct {
}

type AfterDel func(keys ...string)

type Cache interface {
	Set(key string, value interface{}) error
	SetEx(key string, value interface{}, expiration time.Duration) error
	Exists(key string) (error, bool)
	Get(key string, dest any) (error, bool)
	Del(key ...string) (error, bool)
	Keys(key string) ([]string, error)
	AddAfterDel(f AfterDel)
	PublishDelKey(key ...string) error
}

// ----------
type cacheAdapter struct {
}

func (c cacheAdapter) Set(key string, value interface{}) error {
	return c.SetEx(key, value, 0)
}

func (c cacheAdapter) SetEx(key string, value interface{}, expiration time.Duration) error {
	logger.Errorf("this is a cacheAdapter, will ignore set")
	return nil
}

func (c cacheAdapter) Exists(key string) (error, bool) {
	logger.Errorf("this is a cacheAdapter, will ignore Get")
	return nil, false
}

func (c cacheAdapter) Get(key string, dest interface{}) (error, bool) {
	return c.GetEx(key, dest, 0)
}

func (c cacheAdapter) GetEx(key string, dest interface{}, expiration time.Duration) (error, bool) {
	logger.Errorf("this is a cacheAdapter, will ignore Get")
	return nil, false
}

func (c cacheAdapter) Del(key ...string) (error, bool) {
	return nil, false
}

func (c cacheAdapter) Keys(key string) ([]string, error) {
	logger.Errorf("this is a cacheAdapter, will ignore Get")
	return nil, nil
}
func (c cacheAdapter) AddAfterDel(f AfterDel) {
}
func (c cacheAdapter) PublishDelKey(key ...string) error {
	logger.Errorf("this is a cacheAdapter, will ignore PublishDelKey")
	return nil
}
