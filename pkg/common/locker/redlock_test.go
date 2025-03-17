package locker_test

import (
	"github.com/davycun/eta/pkg/common/global"
	_ "github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/locker"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestLock(t *testing.T) {

	lk := locker.NewRedLock([]*redis.Client{global.GetRedis()}, func(ro *locker.RedLockOption) {
		ro.TTL = time.Second * 4
	})
	assert.True(t, lk.Acquire())
	w := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			l := locker.NewRedLock([]*redis.Client{global.GetRedis()})
			assert.False(t, l.Acquire())
		}()
	}
	w.Wait()

	//判断是否会自动续期
	time.Sleep(lk.TTL + time.Second*2)
	for i := 0; i < 20; i++ {
		w.Add(1)
		go func() {
			defer w.Done()
			x := i
			l := locker.NewRedLock([]*redis.Client{global.GetRedis()})
			rs := l.Acquire()
			logger.Infof("i=%d, 锁{lockKey:%s,lockValue:%s}获取结果:%t", x, l.LockKey, l.LockValue, rs)
			assert.False(t, rs)
		}()
	}
	w.Wait()
	lk.Release()

	l := locker.NewRedLock([]*redis.Client{global.GetRedis()})
	assert.True(t, l.Acquire())
	l.Release()

}
