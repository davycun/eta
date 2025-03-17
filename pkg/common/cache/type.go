package cache

import (
	"errors"
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
	Exists(key string) (bool, error)
	Get(key string, dest any) (bool, error)
	Del(key ...string) (bool, error)
	TTL(key string) (time.Duration, error)
	Keys(key string) ([]string, error)
	AddAfterDel(f AfterDel)
	PublishDelKey(key ...string) error
	Detail(key string) (any, time.Duration, error)
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

func (c cacheAdapter) Exists(key string) (bool, error) {
	logger.Errorf("this is a cacheAdapter, will ignore Get")
	return false, nil
}

func (c cacheAdapter) Get(key string, dest interface{}) (bool, error) {
	logger.Errorf("this is a cacheAdapter, will ignore Get")
	return false, nil
}

func (c cacheAdapter) Del(key ...string) (bool, error) {
	return false, nil
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
func (c cacheAdapter) Detail(key string) (any, time.Duration, error) {
	return nil, 0, errors.New("this is a cacheAdapter, will ignore Detail")
}
func (c cacheAdapter) TTL(key string) (time.Duration, error) {
	return 0, errors.New("this is a cacheAdapter, will ignore TTL")
}
