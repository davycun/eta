package cache

import (
	"context"
	"errors"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

const (
	deleteKeyChannel = "deleteKeyChannel"
)

type redisCache struct {
	cacheAdapter
	rds       *redis.Client
	local     sync.Map
	expireKey sync.Map
	afterDel  []AfterDel
}

func (c *redisCache) Set(key string, value interface{}) error {
	var (
		marshal []byte
		err     error
	)
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if value != nil {
				marshal, err = jsoniter.Marshal(value)
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			_, err = Del(key)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			return c.rds.Set(context.Background(), key, marshal, 0).Err()
		}).Err
}
func (c *redisCache) SetEx(key string, value interface{}, expiration time.Duration) error {
	var (
		marshal []byte
		err     error
	)
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if value != nil {
				marshal, err = jsoniter.Marshal(value)
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			_, err = Del(key)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if expiration > 0 {
				expire := time.Now().Add(expiration)
				c.expireKey.Store(key, expire)
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			return c.rds.SetEx(context.Background(), key, marshal, expiration).Err()
		}).Err
}
func (c *redisCache) Get(key string, dest any) (bool, error) {
	var (
		err       error
		bs        []byte
		value, ok = c.local.Load(key)
	)
	if ok {
		bs = value.([]byte)
	} else {
		val := c.rds.Get(context.Background(), key)
		if errors.Is(val.Err(), redis.Nil) {
			return false, nil
		}
		bs, err = val.Bytes()
		if err != nil {
			return false, err
		}

		ttl := c.rds.TTL(context.Background(), key)
		dur := ttl.Val()
		if dur > 0 {
			//本地缓存提前30秒到期
			expire := time.Now().Add(dur - (time.Second * 30))
			c.expireKey.Store(key, expire)
		}
		c.local.Store(key, bs)
	}

	if len(bs) > 0 {
		err = jsoniter.Unmarshal(bs, dest)
		if err != nil {
			return false, err
		}
	}

	return true, err
}
func (c *redisCache) Detail(key string) (any, time.Duration, error) {
	var (
		err    error
		dur    time.Duration
		strCmd *redis.StringCmd
		bs     []byte
		val    any
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			strCmd = c.rds.Get(context.TODO(), key)
			return strCmd.Err()
		}).
		Call(func(cl *caller.Caller) error {
			bs, err = strCmd.Bytes()
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if len(bs) > 0 {
				return jsoniter.Unmarshal(bs, &val)
			}
			return nil
		}).Err

	if errors.Is(err, redis.Nil) {
		return nil, 0, nil
	}

	ttl := c.rds.TTL(context.TODO(), key)
	dur, err = ttl.Result()

	return val, dur, err
}
func (c *redisCache) Exists(key string) (bool, error) {
	_, ok := c.local.Load(key)
	if ok {
		return true, nil
	} else {
		val := c.rds.Exists(context.Background(), key)
		if errors.Is(val.Err(), redis.Nil) || val.Val() == 0 {
			return false, nil
		}
		return true, nil
	}
}
func (c *redisCache) TTL(key string) (time.Duration, error) {
	val := c.rds.TTL(context.Background(), key)
	if errors.Is(val.Err(), redis.Nil) {
		return 0, nil
	}
	return val.Val(), val.Err()
}
func (c *redisCache) Expire(key string, expiration time.Duration) (error, bool) {
	val := c.rds.Expire(context.Background(), key, expiration)
	return val.Err(), val.Val()
}
func (c *redisCache) Del(keys ...string) (bool, error) {
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			del := c.rds.Del(context.Background(), keys...)
			if del.Err() != nil && errors.Is(del.Err(), redis.Nil) {
				return del.Err()
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			for _, v := range keys {
				c.local.Delete(v)
				c.expireKey.Delete(v)

				//删除缓存后，通知所有的redis 客户端本地缓存进行删除
				pub := c.rds.Publish(context.Background(), deleteKeyChannel, v)
				if pub.Err() != nil && errors.Is(pub.Err(), redis.Nil) {
					return pub.Err()
				}
			}
			return nil
		}).Err

	return err == nil, err
}

func (c *redisCache) Keys(key string) (keys []string, err error) {
	keysCmd := c.rds.Keys(context.Background(), key)
	keys, err = keysCmd.Result()
	return
}
func (c *redisCache) Scan(cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error) {
	cmd := c.rds.Scan(context.Background(), cursor, match, count)
	keys, newCursor, err = cmd.Result()
	return
}
func (c *redisCache) Unlink(keys ...string) error {
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return c.rds.Unlink(context.TODO(), keys...).Err()
		}).
		Call(func(cl *caller.Caller) error {
			for _, v := range keys {
				c.local.Delete(v)
				c.expireKey.Delete(v)
				//删除缓存后，通知所有的redis 客户端本地缓存进行删除
				return c.rds.Publish(context.TODO(), deleteKeyChannel, v).Err()
			}
			return nil
		}).Err
}

func (c *redisCache) AddAfterDel(f AfterDel) {
	c.afterDel = append(c.afterDel, f)
}
func (c *redisCache) callAfterDel(key ...string) {
	for i, _ := range c.afterDel {
		c.afterDel[i](key...)
	}
}
func (c *redisCache) PublishDelKey(key ...string) error {
	for _, v := range key {
		//删除缓存后，通知所有的redis 客户端本地缓存进行删除
		pub := c.rds.Publish(context.Background(), deleteKeyChannel, v)
		if pub.Err() != nil && !errors.Is(pub.Err(), redis.Nil) {
			return pub.Err()
		}
	}
	return nil
}

// 订阅清理本地缓存
func (c *redisCache) subscribeDeleteCommand() {
	var (
		pubSub    = c.rds.Subscribe(context.Background(), deleteKeyChannel)
		rcChannel = pubSub.Channel()
		ticker    = time.NewTicker(time.Minute * 1)
	)

	defer func() {
		err := pubSub.Close()
		if err != nil {
			logger.Errorf("close redis pubSub err %s", err)
		}
		ticker.Stop()
		if r := recover(); r != nil {
			logger.Errorf("cache of redis subscribe from deleteKeyChannel panic %v", r)
		}
	}()

	for {
		select {
		case <-ticker.C:
			nw := time.Now()
			delKeys := make([]string, 0, 10)
			c.expireKey.Range(func(key, value any) bool {
				t, ok := value.(time.Time)
				if ok && nw.After(t) {
					delKeys = append(delKeys, key.(string))
				}
				return true
			})
			if len(delKeys) > 0 {
				logger.Infof("clean the expired key in the expireKey %v,now", delKeys)
			}
			for _, v := range delKeys {
				c.expireKey.Delete(v)
				c.local.Delete(v)
				c.callAfterDel(v)
			}

		case msg := <-rcChannel:
			if msg == nil {
				continue
			}
			logger.Infof("receive msg %s from %s", msg.Payload, msg.Channel)
			c.local.Delete(msg.Payload)
			c.callAfterDel(msg.Payload)
		}
	}
}
