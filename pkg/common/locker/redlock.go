package locker

import (
	"context"
	"errors"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

// RedLock
// 算法的核心思想是：客户端尝试在一组独立的 Redis 节点上获取锁，只有在大多数节点（N/2+1）上成功获取锁，并且锁的获取时间在合理的范围内（即未超过锁的失效时间），才认为锁获取成功
// 注意这个不可以夸协程并发调用，最好用一次New一次
type RedLock struct {
	redisClient   []*redis.Client // redis客户端
	retryCount    int             //获取锁尝试的次数
	ttl           time.Duration   //锁的失效时间，time to live
	lockKey       string
	lockValue     string
	wg            *sync.WaitGroup //锁的等待组
	successClient []*redis.Client //为了续期用
	stopRenewal   chan struct{}   //停止续期的信号
	renewalWg     *sync.WaitGroup //等待续期协程结束
}

func (l *RedLock) Acquire() (bool, error) {

	if len(l.successClient) > 0 {
		return false, errors.New("已经获取过锁，无需再次获取")
	}

	var (
		wg           = &sync.WaitGroup{}
		c            = context.Background()
		failedClient = make([]string, 0)
		errList      = make([]error, 0)
		startTime    = time.Now() //为了计算总的获取锁的时间不超过超时时间，避免有一部分成功获取锁的redis已经超时了
	)
	for _, rc := range l.redisClient {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b, err := rc.SetNX(c, l.lockKey, l.lockValue, l.ttl).Result()
			if err != nil {
				errList = append(errList, err)
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
	if (len(l.redisClient)/2+1) <= len(failedClient) || endTime.Sub(startTime) >= l.ttl {
		_, _ = l.Release() //如果失败，记得释放成功的锁
		return false, errors.Join(errList...)
	}
	l.startRenewal(c)
	return true, nil
}

func (l *RedLock) Release() (bool, error) {

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

	return true, nil
}

// 启动锁续期
func (l *RedLock) startRenewal(c context.Context) {
	var (
		tl = time.NewTicker(l.ttl / 3)
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
	_, err := rc.Eval(c, releaseScript, []string{l.lockKey}, l.lockValue).Result()
	if err != nil {
		logger.Errorf("Failed to release lock on %s: %s", rc.String(), err)
	}
}

func (l *RedLock) renewOne(c context.Context, rc *redis.Client) {
	var (
		renewScript = `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("pexpire", KEYS[1], ARGV[2])
		else
			return 0
		end`
	)
	_, err := rc.Eval(c, renewScript, []string{l.lockKey}, l.lockValue, int64(l.ttl.Seconds())).Result()
	if err != nil {
		logger.Errorf("Failed to renew lock on %s: %s", rc.String(), err)
	}
}
