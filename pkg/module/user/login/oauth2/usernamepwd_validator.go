package oauth2

import (
	"context"
	"fmt"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/setting"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// LoginFailLockCheck 登录失败锁定-检查是否锁定
func LoginFailLockCheck(db *gorm.DB, username string) (err error) {
	lc, _ := setting.GetLoginConfig(db)
	if !lc.FailLock() {
		return
	}
	lockKey := constants.RedisKey(constants.LoginFailLock, username)
	rds := global.GetRedis()
	ttl := rds.TTL(context.Background(), lockKey)
	dur := ttl.Val()
	if dur.Seconds() > 0 {
		err = errs.NewClientError(fmt.Sprintf("登录验证失败次数过多,请%s后再试！", utils.FmtDuration(dur, "小时", "分钟")))
	}
	return
}

// LoginFailLockCounterIncr 登录失败锁定-登录失败时，计数加1并且判断是否执行锁定
func LoginFailLockCounterIncr(db *gorm.DB, username string) (err error) {
	lc, _ := setting.GetLoginConfig(db)
	if !lc.FailLock() {
		return
	}
	counterKeys := make([]string, lc.FailLockDurationMinutes)
	dtNow := time.Now()
	for i := int64(0); i < lc.FailLockDurationMinutes; i++ {
		dt := dtNow.Add(-time.Minute * time.Duration(i))
		counterKeys[i] = constants.RedisKey(constants.LoginFailCounter, username, dt.Format("200601021504"))
	}
	rds := global.GetRedis()
	r := rds.IncrBy(context.Background(), counterKeys[0], 1)
	if r.Err() != nil {
		err = r.Err()
		return
	}
	if r.Val() == 1 {
		err = rds.Expire(context.Background(), counterKeys[0], time.Minute*time.Duration(lc.FailLockDurationMinutes)).Err()
	}
	counts := rds.MGet(context.Background(), counterKeys...)
	if counts.Err() != nil {
		err = counts.Err()
		return
	}
	total := 0
	for _, v := range counts.Val() {
		if v != nil {
			num, err1 := strconv.Atoi(v.(string))
			if err1 == nil {
				total += num
			}
		}
	}
	if int64(total) >= lc.FailLockMaxTimes {
		err = LoginFailLockLock(db, username)
		if err != nil {
			return
		}
	}

	return
}

// LoginFailLockLock 登录失败锁定-锁定
func LoginFailLockLock(db *gorm.DB, username string) (err error) {
	lc, _ := setting.GetLoginConfig(db)
	if !lc.FailLock() {
		return
	}
	lockKey := constants.RedisKey(constants.LoginFailLock, username)
	rds := global.GetRedis()
	return rds.SetEx(context.Background(), lockKey, "1", time.Minute*time.Duration(lc.FailLockLockMinutes)).Err()
}

// LoginFailLockCounterClear 登录失败锁定-清除计数
func LoginFailLockCounterClear(db *gorm.DB, username string) (err error) {
	lc, _ := setting.GetLoginConfig(db)
	if !lc.FailLock() {
		return
	}

	toDelKeys := make([]string, lc.FailLockDurationMinutes)
	dtNow := time.Now()
	for i := int64(0); i < lc.FailLockDurationMinutes; i++ {
		dt := dtNow.Add(-time.Minute * time.Duration(i))
		toDelKeys[i] = constants.RedisKey(constants.LoginFailCounter, username, dt.Format("200601021504"))
	}
	rds := global.GetRedis()
	return rds.Del(context.Background(), toDelKeys...).Err()
}
