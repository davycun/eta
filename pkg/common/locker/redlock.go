package locker

import (
	"context"
	"github.com/davycun/eta/pkg/common/id/nanoid"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

// RedLock
// 算法的核心思想是：客户端尝试在一组独立的 Redis 节点上获取锁，只有在大多数节点（N/2+1）上成功获取锁，并且锁的获取时间在合理的范围内（即未超过锁的失效时间），才认为锁获取成功
// 注意这个不可以夸协程并发调用，最好用一次New一次
type RedLockOption struct {
	RetryCount int           //获取锁尝试的次数
	TTL        time.Duration //锁的失效时间，time to live，NewRedLock默认30秒
	LockKey    string
	LockValue  string
}

type SetOptions func(ro *RedLockOption)
type RedLock struct {
	RedLockOption
	redisClient   []*redis.Client // redis客户端
	successClient []*redis.Client //为了续期用
	stopRenewal   chan struct{}   //停止续期的信号
	renewalWg     *sync.WaitGroup //等待续期协程结束
}

// NewRedLock
// 注意这里不同的业务一定要自己指定LockKey，否则都用默认的key，可能会发生冲突
func NewRedLock(redisClient []*redis.Client, options ...SetOptions) *RedLock {
	rl := &RedLock{
		redisClient:   redisClient,
		renewalWg:     &sync.WaitGroup{},
		successClient: make([]*redis.Client, 0),
		stopRenewal:   make(chan struct{}),
	}
	for _, fc := range options {
		fc(&rl.RedLockOption)
	}
	if rl.RetryCount < 1 {
		rl.RetryCount = 2
	}
	if rl.LockKey == "" {
		rl.LockKey = rl.randKey()
	}
	if rl.LockValue == "" {
		rl.LockValue = rl.randValue()
	}
	if rl.TTL < time.Millisecond {
		rl.TTL = time.Second * 20 //默认30秒
	}
	return rl
}

func (l *RedLock) randKey() string {
	return "eta:lock:redlock"
}
func (l *RedLock) randValue() string {
	return nanoid.New()
}

// Acquire
// TODO 支持可重入
func (l *RedLock) Acquire() bool {
	if len(l.successClient) > 0 {
		logger.Warn("已经获取过锁，无需再次获取")
		return false
	}

	var (
		wg           = &sync.WaitGroup{}
		c            = context.Background()
		failedClient = make([]string, 0)
		startTime    = time.Now() //为了计算总的获取锁的时间不超过超时时间，避免有一部分成功获取锁的redis已经超时了
	)
	for _, rc := range l.redisClient {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b, err := rc.SetNX(c, l.LockKey, l.LockValue, l.TTL).Result()
			if err != nil {
				logger.Errorf("locked err with redis %s : %s", rc.String(), err)
			}
			if !b {
				failedClient = append(failedClient, rc.String())
			} else {
				l.successClient = append(l.successClient, rc)
			}
		}()
	}
	//等待所有 goroutine 执行完成
	wg.Wait()

	//超过半数据节点获取锁失败，则认为锁获取失败
	//如果获取锁的总时间超过ttl，则认为锁获取失败，如果超过ttl了代表之前有一部分获取成功的锁可能已经超时，在redis服务端已经过期了
	endTime := time.Now()
	if (len(l.redisClient)/2 + 1) <= len(failedClient) {
		l.Release() //如果失败，记得释放成功的锁
		logger.Warnf("获取锁失败，客户端数量[%d]，失败数量[%d]", len(l.redisClient), len(failedClient))
		return false
	}
	tm := endTime.Sub(startTime)
	if tm >= l.TTL {
		logger.Warnf("获取锁失败，因为超时，超时阈值[%d]秒，实际时间[%d]秒", l.TTL/time.Second, tm/time.Second)
		l.Release()
		return false
	}
	l.startRenewal(c)
	return true
}

func (l *RedLock) Release() bool {

	if len(l.successClient) < 1 {
		return true
	}

	var (
		c  = context.Background()
		wg = &sync.WaitGroup{}
	)
	//停止续期
	close(l.stopRenewal)
	//等待需求协程结束
	l.renewalWg.Wait()
	for _, rc := range l.successClient {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.releaseOne(c, rc)
		}()
	}
	wg.Wait()
	return true
}

// 启动锁续期
func (l *RedLock) startRenewal(c context.Context) {
	var (
		tl = time.NewTicker(l.TTL / 2)
	)
	l.renewalWg.Add(1)
	go func() {
		defer l.renewalWg.Done()
		defer tl.Stop()
		for {
			select {
			case <-tl.C:
				w := &sync.WaitGroup{}
				for _, rc := range l.successClient {
					w.Add(1)
					go func() {
						defer w.Done()
						l.renewOne(c, rc)
					}()
				}
				w.Wait()
			case <-l.stopRenewal:
				//停止续期
				logger.Infof("{lockKey:%s,lockValue:%s} 续期已经停止", l.LockKey, l.LockValue)
				return
			}
		}
	}()
}

func (l *RedLock) releaseOne(c context.Context, rc *redis.Client) {
	var (
		releaseScript = `
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end`
	)
	v, err := rc.Eval(c, releaseScript, []string{l.LockKey}, l.LockValue).Result()
	logger.Infof("正在给redis %s redlock {lockKey:%s,lockValue:%s} 释放锁，释放的结果%v", rc.String(), l.LockKey, l.LockValue, v)
	if err != nil {
		logger.Errorf("Failed to release lock on %s: %s", rc.String(), err)
	}
}

func (l *RedLock) renewOne(c context.Context, rc *redis.Client) {
	var (
		renewScript = `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("PEXPIRE", KEYS[1], ARGV[2])
		else
			return 0
		end`
	)
	v, err := rc.Eval(c, renewScript, []string{l.LockKey}, l.LockValue, l.TTL.Milliseconds()).Result()
	logger.Infof("正在给redis %s redlock {lockKey:%s,lockValue:%s} 续期，续期的结果%v", rc.String(), l.LockKey, l.LockValue, v)
	if err != nil {
		logger.Errorf("Failed to renew lock on %s: %s", rc.String(), err)
	}
}
